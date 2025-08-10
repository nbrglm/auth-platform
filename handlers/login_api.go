package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/internal"
	"github.com/nbrglm/auth-platform/internal/auth"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type LoginAPIHandler struct {
	LoginCounter *prometheus.CounterVec
}

func NewLoginAPIHandler() *LoginAPIHandler {
	return &LoginAPIHandler{
		LoginCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_login_api_requests",
				Help:      "Total number of user login API requests",
			},
			[]string{"status"},
		),
	}
}

func (h *LoginAPIHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.LoginCounter)
	engine.POST("/api/auth/login", middlewares.RateLimitAPIMiddleware(), h.HandleLoginAPI)
}

// HandleLoginAPI godoc
// @Summary User Login API
// @Description Handles user login requests via API.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.UserLoginData true "User Login Data"
// @Success 200 {object} auth.UserLoginResult "User Login Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/login [post]
func (h *LoginAPIHandler) HandleLoginAPI(c *gin.Context) {
	h.LoginCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "login_api")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var loginData auth.UserLoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		h.LoginCounter.WithLabelValues("invalid_input").Inc()
		log.Debug("Failed to bind login data", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid input data", "Bad Request", http.StatusBadRequest, err).Filter())
		return
	}

	// call the login handler
	result, err := auth.HandleLogin(ctx, log, loginData)
	result = ProcessAPIResult(c, result, err, span, log, h.LoginCounter, "login_api")

	if result == nil {
		return
	}

	// We do not handle special cases like we do earlier in the web handler:
	// - We do not redirect to a verification page if the email is not verified, BUT
	//   we do send a verification email if the email is not verified and return that the email has been sent.
	// - We do not handle multitenant flows here, as this is an API endpoint
	//   and we expect the client to handle the flow continuation (with the flowId returned from here, using flow endpoints like org-select or mfa later).
	// - We do not set tokens in cookies, as this is an API endpoint.
	//   Instead, we return the tokens in the response.
	//   The client should handle storing these tokens (e.g., in local storage or cookies).
	//   This is a design choice to keep the API stateless and allow clients to handle
	//   authentication and session management as they see fit.

	h.LoginCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, result)
}
