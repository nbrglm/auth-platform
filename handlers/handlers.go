package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nbrglm/nexeres/handlers/admin_handlers"
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
		admin_handlers.NewAdminLoginHandler(),
		admin_handlers.NewConfigHandler(),
	}

	// Register API routes
	for _, handler := range handlers {
		handler.Register(engine)
	}

	// NOTE: For resetting the inactivity timer on admin routes,
	// we need to call `middlewares.AdminInactivityReset()` after all computations
	// in each handler are done.
}
