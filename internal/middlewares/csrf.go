package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/internal/logging"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/opts"
	"go.uber.org/zap"
)

// CSRFMiddleware returns a middleware that protects against CSRF attacks.
// It uses the Gorilla CSRF library to validate CSRF tokens in requests.
func CSRFMiddleware() gin.HandlerFunc {
	var csrfMiddleware func(http.Handler) http.Handler = csrf.Protect(
		[]byte(config.Security.CSRF.Secret),
		csrf.Secure(!opts.Debug),          // Use secure cookies in production
		csrf.MaxAge(0),                    // Set session-only CSRF token, no expiration
		csrf.Domain(config.Public.Domain), // Set the domain for the CSRF cookie
		csrf.CookieName(config.Security.CSRF.TokenName),
		csrf.FieldName(config.Security.CSRF.TokenName),
		csrf.RequestHeader(config.Security.CSRF.TokenHeader),
		csrf.HttpOnly(true),                    // Make the CSRF cookie HTTP-only
		csrf.SameSite(csrf.SameSiteStrictMode), // Use Strict SameSite policy for CSRF cookie
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			body, err := json.Marshal(models.NewErrorResponse("This request was blocked to protect you from potential security threats. Please try again or contact support if the issue persists.", "CSRF token validation failed! Please ensure that your request includes a valid CSRF token.", http.StatusForbidden, nil).Filter())
			if err != nil {
				logging.Logger.Error("Failed to marshal CSRF error response", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
				return
			}
			w.WriteHeader(http.StatusForbidden)
			w.Write(body)
		})),
	)

	return adapter.Wrap(csrfMiddleware)

}
