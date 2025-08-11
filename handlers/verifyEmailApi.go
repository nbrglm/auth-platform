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
)

type VerifyEmailAPIHandler struct {
	SendEmailCounter   *prometheus.CounterVec
	VerifyTokenCounter *prometheus.CounterVec
}

func NewVerifyEmailAPIHandler() *VerifyEmailAPIHandler {
	return &VerifyEmailAPIHandler{
		SendEmailCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "verify_email_send_api_requests",
				Help:      "Total number of verify email send API requests",
			},
			[]string{"status"},
		),
		VerifyTokenCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "verify_email_token_api_requests",
				Help:      "Total number of verify email token API requests",
			},
			[]string{"status"},
		),
	}
}

func (h *VerifyEmailAPIHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.SendEmailCounter)
	metrics.Collectors = append(metrics.Collectors, h.VerifyTokenCounter)

	engine.POST("/api/auth/verify-email/send", middlewares.RateLimitAPIMiddleware(), h.HandleSendVerificationEmailAPI)
	engine.POST("/api/auth/verify-email/verify", middlewares.RateLimitAPIMiddleware(), h.HandleVerifyEmailTokenAPI)
}

// HandleSendVerificationEmailAPI godoc
// @Summary Send Verification Email API
// @Description Sends a verification email to the user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.SendVerificationEmailData true "Send Verification Email Data"
// @Success 200 {object} auth.SendVerificationEmailResult "Send Verification Email Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/verify-email/send [post]
func (h *VerifyEmailAPIHandler) HandleSendVerificationEmailAPI(c *gin.Context) {
	h.SendEmailCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "send_verification_email_api")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var emailData auth.SendVerificationEmailData
	if err := c.ShouldBindJSON(&emailData); err != nil {
		ProcessAPIResult[any](c, nil, models.NewErrorResponse("Invalid input data", "Failed to bind JSON data", http.StatusBadRequest, nil), span, log, h.SendEmailCounter, "send_verification_email_api")
		return
	}

	result, err := auth.SendVerificationEmail(ctx, log, emailData)
	result = ProcessAPIResult(c, result, err, span, log, h.SendEmailCounter, "send_verification_email_api")

	if result == nil {
		return
	}

	h.SendEmailCounter.WithLabelValues("success").Inc()

	c.JSON(http.StatusOK, result)
}

// HandleVerifyEmailTokenAPI godoc
// @Summary Verify Email Token API
// @Description Verifies the email using the provided token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.VerifyEmailTokenData true "Verify Email Token Data"
// @Success 200 {object} auth.VerifyEmailTokenResult "Verify Email Token Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/verify-email/verify [post]
func (h *VerifyEmailAPIHandler) HandleVerifyEmailTokenAPI(c *gin.Context) {
	h.VerifyTokenCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "verify_email_token_api")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var tokenData auth.VerifyEmailTokenData
	if err := c.ShouldBindJSON(&tokenData); err != nil {
		ProcessAPIResult[any](c, nil, models.NewErrorResponse("Invalid input data", "Failed to bind JSON data", http.StatusBadRequest, nil), span, log, h.VerifyTokenCounter, "verify_email_token_api")
		return
	}

	result, err := auth.VerifyEmailToken(ctx, log, tokenData)
	result = ProcessAPIResult(c, result, err, span, log, h.VerifyTokenCounter, "verify_email_token_api")

	if result == nil {
		return
	}

	h.VerifyTokenCounter.WithLabelValues("success").Inc()

	c.JSON(http.StatusOK, result)
}
