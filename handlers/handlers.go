package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler interface {
	Register(engine *gin.Engine)
}

// RegisterAPIRoutes registers all API routes for the application.
func RegisterAPIRoutes(engine *gin.Engine) {
	// Register API routes
	signupAPIHandler := NewSignupAPIHandler()
	signupAPIHandler.Register(engine)

	verifyEmailAPIHandler := NewVerifyEmailAPIHandler()
	verifyEmailAPIHandler.Register(engine)

	loginAPIHandler := NewLoginAPIHandler()
	loginAPIHandler.Register(engine)

	// Not API, but not UI either, usually comes under /auth/s/* (s stands for sensitive, meaning these routes require sensitive data like refresh token access)
	refreshTokenHandler := NewRefreshTokenHandler()
	refreshTokenHandler.Register(engine)

	logoutHandler := NewLogoutHandler()
	logoutHandler.Register(engine)

	// UI routes
	signupWEBHandler := NewSignupWEBHandler()
	signupWEBHandler.Register(engine)

	verifyEmailWEBHandler := NewVerifyEmailWEBHandler()
	verifyEmailWEBHandler.Register(engine)

	loginWEBHandler := NewLoginWEBHandler()
	loginWEBHandler.Register(engine)

	unsafeRedirectHandler := NewUnsafeRedirectHandler()
	unsafeRedirectHandler.Register(engine)
}

// ProcessAPIResult is a utility function to handle errors in a consistent way across handlers.
// It logs the error, increments the appropriate Prometheus counter, and sends a JSON response to the client.
// Returns true if an error occurred and was handled, false otherwise.
func ProcessAPIResult[T any](c *gin.Context, response *T, err *models.ErrorResponse, span trace.Span, log *zap.Logger, counter *prometheus.CounterVec, opName string) *T {
	if response != nil {
		return response
	}

	// If response is nil, we handle the error
	counter.WithLabelValues("error").Inc()
	log.Debug("Error occurred during operation!", zap.String("operation", opName), zap.Error(err))

	if err == nil {
		// If there is no error, we return nil to indicate that the operation was not
		err = models.NewGenericErrorResponse(opName, fmt.Errorf("response is unexpectedly null for operation: %s", opName))
	}

	if err.UnderlyingError != nil {
		// Log and Record the underlying error if it exists
		log.Error("Failed to handle operation", zap.String("operation", opName), zap.Error(err.UnderlyingError))
		span.RecordError(err.UnderlyingError)
	}
	span.SetStatus(codes.Error, err.DebugMessage)
	c.JSON(err.Code, err.Filter())

	return nil
}

// ProcessUiResult is a utility function to handle errors in a consistent way across UI handlers.
// It executes the provided function and handles any errors that occur.
// It logs the error, increments the appropriate Prometheus counter, and redirects the user to a specified URL.
// It also stores error messages and return information in the session.
// Returns the result of the provided function if successful, or nil if an error occurred.
func ProcessUiResult[T any](c *gin.Context, response *T, err *models.ErrorResponse, span trace.Span, log *zap.Logger, counter *prometheus.CounterVec, opName string) *T {
	if response != nil {
		// If response is not nil, we assume the operation was successful
		return response
	}

	// If response is nil, we handle the error
	counter.WithLabelValues("error").Inc()
	log.Debug("Error occurred during operation!", zap.String("operation", opName), zap.Error(err))

	if err == nil {
		// If there is no error, we return nil to indicate that the operation was not successful as response is still nil
		err = models.NewGenericErrorResponse(opName, fmt.Errorf("response is unexpectedly null for operation: %s", opName))
	}

	if err.UnderlyingError != nil {
		// Log and Record the underlying error if it exists
		log.Error("Failed to handle operation", zap.String("operation", opName), zap.Error(err.UnderlyingError))
		span.RecordError(err.UnderlyingError)
	}
	span.SetStatus(codes.Error, err.DebugMessage)
	setErrorAndRedirect(c, err)

	// If we reach here, it means we have handled the error and redirected the user.
	// We return nil to indicate that the operation was not successful.
	return nil
}

func setErrorAndRedirect(c *gin.Context, errResponse *models.ErrorResponse) {
	middlewares.SetPageError(c, errResponse.Message, errResponse.UIParams.RetryButtonText, errResponse.UIParams.RetryUrl)
	c.Redirect(http.StatusSeeOther, errResponse.UIParams.RedirectUrl)
}

// getPageError is a utility function to retrieve the page error from the context.
// The title is used to set the title of the error page.
func getPageError(c *gin.Context, title string) *models.CommonPageParams {
	errorMessage := c.GetString(middlewares.CtxPageErrorKey)
	returnTo := c.GetString(middlewares.CtxPageErrorReturnURLKey)
	returnButtonText := c.GetString(middlewares.CtxPageErrorReturnButtonTextKey)

	if errorMessage == "" {
		errorMessage = "An unknown error occurred. Please try again."
	}

	pageErr := models.NewPageError(title, errorMessage, returnTo, returnButtonText)

	return &pageErr
}
