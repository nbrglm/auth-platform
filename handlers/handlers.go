package handlers

import (
	"github.com/gin-gonic/gin"
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
	handlers := []Handler{
		NewSignupHandler(),
		NewVerifyEmailHandler(),
		NewLoginHandler(),
		NewRefreshTokenHandler(),
		NewLogoutHandler(),
	}

	// Register API routes
	for _, handler := range handlers {
		handler.Register(engine)
	}
}

// ProcessError is a utility function to handle errors in a consistent way across handlers.
// It logs the error, increments the appropriate Prometheus counter, and sends a JSON response to the client.
// Returns true if an error occurred and was handled, false otherwise.
func ProcessError(c *gin.Context, err *models.ErrorResponse, span trace.Span, log *zap.Logger, counter *prometheus.CounterVec, opName string) bool {
	if err == nil {
		return false
	}

	counter.WithLabelValues("error").Inc()
	log.Debug("Error occurred during operation!", zap.String("operation", opName), zap.Error(err))

	if err.UnderlyingError != nil {
		// Log and Record the underlying error if it exists
		log.Error("Failed to handle operation", zap.String("operation", opName), zap.Error(err.UnderlyingError))
		span.RecordError(err.UnderlyingError)
	}
	span.SetStatus(codes.Error, err.DebugMessage)
	c.JSON(err.Code, err.Filter())
	return true
}
