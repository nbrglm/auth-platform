package admin_handlers

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nbrglm/nexeres/config"
	"github.com/nbrglm/nexeres/internal"
	"github.com/nbrglm/nexeres/internal/cache"
	"github.com/nbrglm/nexeres/internal/metrics"
	"github.com/nbrglm/nexeres/internal/models"
	"github.com/nbrglm/nexeres/internal/notifications"
	"github.com/nbrglm/nexeres/internal/otp"
	"github.com/nbrglm/nexeres/internal/tokens"
	"github.com/nbrglm/nexeres/opts"
	"github.com/nbrglm/nexeres/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type AdminLoginHandler struct {
	AdminLoginCounter       *prometheus.CounterVec
	AdminLoginVerifyCounter *prometheus.CounterVec
}

func NewAdminLoginHandler() *AdminLoginHandler {
	return &AdminLoginHandler{
		AdminLoginCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "admin",
				Name:      "login_requests_total",
				Help:      "Total number of admin login requests",
			},
			[]string{"status"},
		),
		AdminLoginVerifyCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "admin",
				Name:      "login_verify_requests_total",
				Help:      "Total number of admin login verification requests",
			},
			[]string{"status"},
		),
	}
}

func (h *AdminLoginHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.AdminLoginCounter, h.AdminLoginVerifyCounter)
	engine.POST("/api/admin/login", h.AdminLogin)
	engine.POST("/api/admin/login/verify", h.VerifyAdminLogin)
}

type AdminLoginData struct {
	Email string `json:"email" binding:"required,email"`
}

type AdminLoginResult struct {
	// If true, the email is an existing admin email and a verification code has been sent.
	// The user can then proceed to verify the code.
	EmailSent bool `json:"emailSent"`

	// The flow ID for the login process. This is required for verifying the code.
	FlowId string `json:"flowId"`
}

// AdminLogin godoc
// @Summary Admin login
// @Description Sends a login code to the admin's email if it exists
// @Tags admin
// @Accept json
// @Produce json
// @Param data body AdminLoginData true "Admin login data"
// @Success 200 {object} AdminLoginResult "Admin Login Result"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/admin/login [post]
func (h *AdminLoginHandler) AdminLogin(c *gin.Context) {
	h.AdminLoginCounter.WithLabelValues("received").Inc()

	// For single-tenant login, we can directly validate the user's credentials
	ctx, log, span := internal.WithContext(c.Request.Context(), "admin_login")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var requestData AdminLoginData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Invalid request data", "The provided request data is invalid", http.StatusBadRequest, err), span, log, h.AdminLoginCounter, "admin_login")
		return
	}

	if slices.Contains(config.Admins.Emails, requestData.Email) {
		// Generate a new login flow
		flowId, err := uuid.NewV7()
		if err != nil {
			utils.ProcessError(c, models.NewErrorResponse("Internal server error", "Failed to generate flow ID", http.StatusInternalServerError, err), span, log, h.AdminLoginCounter, "admin_login")
			return
		}

		now := time.Now()
		expiresAt := now.Add(time.Duration(config.Admins.SessionTimeout * int(time.Second)))

		code, err := otp.NewAlphaNumericOTP(opts.AdminOTPLength)
		if err != nil {
			utils.ProcessError(c, models.NewErrorResponse("Internal server error", "Failed to generate OTP", http.StatusInternalServerError, err), span, log, h.AdminLoginCounter, "admin_login")
			return
		}

		err = cache.StoreAdminLoginFlow(ctx, cache.AdminLoginFlowData{
			ID:        flowId.String(),
			Email:     requestData.Email,
			Code:      code,
			CreatedAt: now,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			utils.ProcessError(c, models.NewErrorResponse("Internal server error", "Failed to store admin login flow", http.StatusInternalServerError, err), span, log, h.AdminLoginCounter, "admin_login")
			return
		}

		err = notifications.SendAdminLoginEmail(ctx, notifications.SendAdminLoginEmailParams{
			Email:     requestData.Email,
			Code:      code,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			utils.ProcessError(c, models.NewErrorResponse("Unable to send admin login email - Internal Server Error", "Failed to send admin login email", http.StatusInternalServerError, err), span, log, h.AdminLoginCounter, "admin_login")
			return
		}

		h.AdminLoginCounter.WithLabelValues("success").Inc()
		c.JSON(http.StatusOK, AdminLoginResult{
			EmailSent: true,
			FlowId:    flowId.String(),
		})
		return
	}

	// For security reasons, we do not reveal whether the email exists or not
	h.AdminLoginCounter.WithLabelValues("invalid_request").Inc()
	utils.ProcessError(c, models.NewErrorResponse("Invalid request", "The provided email does not exist", http.StatusUnauthorized, nil), span, log, h.AdminLoginCounter, "admin_login")
}

type VerifyAdminLoginData struct {
	FlowId string `json:"flowId" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type VerifyAdminLoginResult struct {
	// If true, the admin login was successful and tokens are provided.
	Success bool `json:"success"`

	// The secure token for the admin session.
	// Expires as configured in the server settings and with idle.
	Token string `json:"token"`
}

// VerifyAdminLogin godoc
// @Summary Verify admin login
// @Description Verifies the admin login using the code sent to email
// @Tags admin
// @Accept json
// @Produce json
// @Param data body VerifyAdminLoginData true "Verify admin login data"
// @Success 200 {string} VerifyAdminLoginResult "Verify Admin Login Result"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/admin/login/verify [post]
func (h *AdminLoginHandler) VerifyAdminLogin(c *gin.Context) {
	h.AdminLoginVerifyCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "verify_admin_login")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var requestData VerifyAdminLoginData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Invalid request data", "The provided request data is invalid", http.StatusBadRequest, err), span, log, h.AdminLoginVerifyCounter, "verify_admin_login")
		return
	}

	flowData, err := cache.GetAdminLoginFlow(ctx, requestData.FlowId)
	if err != nil {
		if err == cache.ErrKeyNotFound {
			utils.ProcessError(c, models.NewErrorResponse("Invalid or expired flow", "The provided flow ID is invalid or has expired", http.StatusUnauthorized, nil), span, log, h.AdminLoginVerifyCounter, "verify_admin_login")
			return
		}
		utils.ProcessError(c, models.NewErrorResponse("Internal server error", "Failed to retrieve admin login flow", http.StatusInternalServerError, err), span, log, h.AdminLoginVerifyCounter, "verify_admin_login")
		return
	}

	if requestData.Code != flowData.Code {
		h.AdminLoginVerifyCounter.WithLabelValues("invalid_request").Inc()
		utils.ProcessError(c, models.NewErrorResponse("Invalid code", "The provided code is incorrect", http.StatusUnauthorized, nil), span, log, h.AdminLoginVerifyCounter, "verify_admin_login")
		return
	}

	// Successful verification, generate token
	token, hash, err := tokens.NewAdminToken()
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Internal server error", "Failed to generate admin token", http.StatusInternalServerError, err), span, log, h.AdminLoginCounter, "verify_admin_login")
		return
	}

	now := time.Now()
	expiresAt := time.Now().Add(time.Duration(config.Admins.SessionTimeout * int(time.Second)))

	err = cache.StoreAdminSession(ctx, cache.AdminSessionData{
		Token:     *hash,
		Email:     flowData.Email,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	})

	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Internal server error", "Failed to store admin session", http.StatusInternalServerError, err), span, log, h.AdminLoginCounter, "verify_admin_login")
		return
	}

	// Delete the used login flow
	err = cache.DeleteAdminLoginFlow(ctx, requestData.FlowId)
	if err != nil {
		if err != cache.ErrKeyNotFound {
			utils.ProcessError(c, models.NewErrorResponse("Internal server error", "Failed to delete used admin login flow", http.StatusInternalServerError, err), span, log, h.AdminLoginCounter, "verify_admin_login")
			return
		}
	}

	c.Header(tokens.AdminTokenExpiryHeaderName, expiresAt.Format(time.RFC3339))

	h.AdminLoginVerifyCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, VerifyAdminLoginResult{
		Success: true,
		Token:   *token,
	})
}
