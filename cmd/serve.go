package cmd

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/handlers"
	"github.com/nbrglm/auth-platform/internal/cache"
	"github.com/nbrglm/auth-platform/internal/logging"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/notifications"
	"github.com/nbrglm/auth-platform/internal/notifications/templates"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"github.com/nbrglm/auth-platform/internal/tracing"
	"github.com/nbrglm/auth-platform/opts"
	"github.com/nbrglm/auth-platform/utils"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.uber.org/zap"
)

func initServeCommand(webFS embed.FS) {
	var serveCommand = &cobra.Command{
		Use:   "serve",
		Short: "Start the server and listen for incoming requests",
		Run: func(cmd *cobra.Command, args []string) {
			// Start the server
			runServer(cmd, webFS)
		},
	}

	serveCommand.Flags().StringVar(opts.ConfigPath, "config", "/etc/nbrglm/workspace/auth-platform/config.yaml", "Path to the config file")
	serveCommand.MarkPersistentFlagFilename("config", "yaml", "yml")

	rootCmd.AddCommand(serveCommand)
}

func runServer(cmd *cobra.Command, webFS embed.FS) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting %s v%s...\n", opts.Name, opts.Version)
	fmt.Printf("Config file: %s\n", *opts.ConfigPath)

	// Initialize the validator before everything else, since validation is used by the config file loader.
	utils.InitValidator()

	// Load the configuration file
	if err := config.LoadConfigOptions(*opts.ConfigPath); err != nil {
		cmd.PrintErrf("Error loading config file: %v\n", err)
		os.Exit(1)
		return
	}

	// Initialize the logger
	logging.InitLogger()

	engine := gin.Default()
	if opts.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Add CORS middleware
	middlewares.InitCORS(engine)

	if !opts.Debug {
		// CSRF if not in debug mode
		engine.Use(middlewares.CSRFMiddleware())
	}

	//TODO: Remove this middleware once the cookie handling is done properly.
	// engine.Use(func(c *gin.Context) {
	// 	fmt.Printf("\nINITIAL: \nRAW COOKIE HEADER: %s\nPATH: %s\n", c.Request.Header.Get("Cookie"), c.Request.URL.Path)
	// 	c.Next()
	// })

	// Add the API Key middleware, before rate limiting middlewares,
	// since those need access to keys for rate limiting
	engine.Use(middlewares.APIKeyMiddleware())

	// Add the returnTo URL middleware
	engine.Use(middlewares.ReturnToURLMiddleware())

	// Add the middleware for storing page errors in the context
	// by retrieving them from the cookies.
	engine.Use(middlewares.PageErrorStorageMiddleware())

	if opts.Debug {
		logging.Logger.Warn("Debug mode is enabled! This is not recommended for production environments. Use with caution. The following behaviour is used.", zap.String("Debug Mode", "Enabled"), zap.String("API Docs", fmt.Sprintf("%s/docs", config.Public.GetBaseURL())), zap.String("CSRF Protection", "Disabled"))
		// Setup docs
		engine.GET("/docs", func(ctx *gin.Context) {
			ctx.Header("Content-Type", "text/html")
			ctx.String(200, `<!doctype html>
	<html>
		<head>
			<title>API Reference</title>
			<meta charset="utf-8" />
			<meta
				name="viewport"
				content="width=device-width, initial-scale=1" />
		</head>
		<body>
			<script
				id="api-reference"
				data-url="/swagger/doc.json"></script>
			<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		</body>
	</html>`)
		})
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Initialize the rate limiter, before adding the handler routes.
	if err := middlewares.InitRateLimitStore(); err != nil {
		logging.Logger.Error("Failed to initialize rate limit store", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	// Register the routes
	handlers.RegisterAPIRoutes(engine)

	// engine.GET("/", func(ctx *gin.Context) {
	// 	claims, exists := ctx.Get(middlewares.CtxSessionTokenClaimsKey)
	// 	if !exists {
	// 		ctx.JSON(http.StatusOK, gin.H{
	// 			"message": "Welcome to the NBRGLM Auth Platform! Please login or signup to continue.",
	// 		})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusOK, claims)
	// })

	// Initialize the metrics collection system
	//
	// NOTE: Always do this after registering the API routes.
	//
	// This is because the collectors need to be registered with the Prometheus registry
	// before the metrics route is added to the engine.
	// And the collectors are only assigned in the Register() methods of each handler, hence we need
	// to call this after registering the API routes.
	// This will also register the /metrics route to serve the metrics in Prometheus format.
	// This is done to ensure that the metrics are collected and reported correctly.
	metrics.InitMetrics()
	metrics.AddMetricsRoute(engine)

	// Initialize the OpenTelemetry Tracer
	err := tracing.InitTracer(context.Background())
	if err != nil {
		// If OTEL tracer provider initialization fails, log the error and exit
		logging.Logger.Error("Failed to initialize OTEL tracer provider", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	tracing.AddTracingMiddleware(engine)

	// Parse the notification templates
	if err := templates.ParseEmailTemplates(); err != nil {
		logging.Logger.Error("Failed to parse email templates", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}
	if err := templates.ParseMessageTemplates(); err != nil {
		logging.Logger.Error("Failed to parse message templates", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	// Setup Notifications senders
	notifications.InitEmail()
	notifications.InitSMS()

	// Initialize the cache
	if err := cache.InitCache(); err != nil {
		logging.Logger.Error("Failed to initialize cache", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	// Initialize the token generation and keys
	if err := tokens.InitTokens(); err != nil {
		logging.Logger.Error("Failed to initialize tokens", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	// Connect with the database
	if err := store.InitDB(context.Background()); err != nil {
		logging.Logger.Error("Failed to initialize database connection pool", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	// Initialize the s3 store
	if err := store.InitS3Store(context.Background()); err != nil {
		logging.Logger.Error("Failed to initialize S3 store", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	// Start the server
	serverAddress := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: engine.Handler(),
	}

	logging.Logger.Info("Starting server", zap.String("address", serverAddress))
	fmt.Printf("Starting server at %v", serverAddress)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Logger.Error("Failed to start server", zap.Error(err))
		logging.ShutdownLogger(context.Background())
		os.Exit(1)
	}

	fmt.Printf("Started server at %v", serverAddress)
	logging.Logger.Info("Started server", zap.String("Address", serverAddress))

	// Wait for OS signals to gracefully shutdown the server
	<-osSignals

	logging.Logger.Info("Received shutdown signal, shutting down server gracefully...")

	logging.Logger.Info("Closing database connection pool")
	if err := store.CloseDB(); err != nil {
		logging.Logger.Error("Failed to close database connection pool", zap.Error(err))
	}

	logging.Logger.Info("Shutting down OTEL tracer provider")
	if err := tracing.ShutdownTracer(context.Background()); err != nil {
		logging.Logger.Error("Failed to shutdown OTEL tracer provider", zap.Error(err))
	}

	// Metrics collector shutdown is not needed as it is handled by the Prometheus registry

	// Shut down server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.Logger.Error("Failed to shutdown server gracefully", zap.Error(err))
	}

	// Wait for the context to be done before exiting
	<-ctx.Done()

	logging.Logger.Info("Shutting down logger.")
	if err := logging.ShutdownLogger(context.Background()); err != nil {
		fmt.Printf("Failed to shutdown logger, %v", err)
	}

	// No logging.* calls after this point, as the logger is shutting down.
}
