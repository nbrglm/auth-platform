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

type SignupAPIHandler struct {
	SignupCounter *prometheus.CounterVec
}

// NewSignupAPIHandler creates a new SignupAPIHandler instance.
func NewSignupAPIHandler() *SignupAPIHandler {
	return &SignupAPIHandler{
		SignupCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_signup_api_requests",
				Help:      "Total number of user signup API requests",
			},
			[]string{"status"},
		),
	}
}

func (h *SignupAPIHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.SignupCounter)
	engine.POST("/api/auth/signup", middlewares.RateLimitAPIMiddleware(), h.HandleSignupAPI)
}

// HandleSignupAPI godoc
// @Summary User Signup API
// @Description Handles user registration requests via API.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.UserSignupData true "User Signup Data"
// @Success 200 {object} auth.UserSignupResult "User Signup Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/signup [post]
func (h *SignupAPIHandler) HandleSignupAPI(c *gin.Context) {
	h.SignupCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "signup_api")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var signupData auth.UserSignupData
	if err := c.ShouldBindJSON(&signupData); err != nil {
		h.SignupCounter.WithLabelValues("invalid_input").Inc()
		log.Debug("Failed to bind signup data", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", "Please check your input and try again.", http.StatusBadRequest, nil).Filter())
		return
	}

	// Call the internal handler
	response, err := auth.HandleSignup(ctx, log, signupData)
	response = ProcessAPIResult(c, response, err, span, log, h.SignupCounter, "SignupAPI")
	if response == nil {
		return // If ProcessAPIResult returns nil, it means an error was handled
	}

	// If we reach here, it means the signup was successful
	h.SignupCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, response)
}
