package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/internal"
	"github.com/nbrglm/auth-platform/internal/auth"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"github.com/prometheus/client_golang/prometheus"
)

type LogoutHandler struct {
	LogoutWEBCounter *prometheus.CounterVec
	LogoutAPICounter *prometheus.CounterVec
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{
		LogoutWEBCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_logout_web_requests",
				Help:      "Total number of user WEB logout requests",
			},
			[]string{"status"},
		),
		LogoutAPICounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_logout_api_requests",
				Help:      "Total number of user API logout requests",
			},
			[]string{"status"},
		),
	}
}

func (h *LogoutHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.LogoutWEBCounter)
	metrics.Collectors = append(metrics.Collectors, h.LogoutAPICounter)

	engine.GET("/auth/s/logout", middlewares.RateLimitUIAuthenticatedMiddleware(), h.HandleLogoutWEB)
	engine.POST("/api/auth/s/logout", middlewares.RateLimitAPIMiddleware(), h.HandleLogoutAPI)
}

func (h *LogoutHandler) HandleLogoutWEB(c *gin.Context) {
	h.LogoutWEBCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "refresh_token_web")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	if !c.GetBool(middlewares.CtxSessionExistsKey) && !c.GetBool(middlewares.CtxSessionRefreshKey) {
		h.LogoutWEBCounter.WithLabelValues("invalid_request").Inc()
		c.Redirect(http.StatusSeeOther, "/auth/login") // Redirect to login page if session is invalid
		return
	}

	// Revoke session
	logoutData := auth.LogoutData{
		SessionToken: c.GetString(middlewares.CtxSessionTokenKey),
		RefreshToken: c.GetString(middlewares.CtxSessionRefreshTokenKey),
	}

	logoutResult, errResp := auth.HandleLogout(ctx, log, logoutData)

	// Remove tokens from cookies regardless of success or failure
	tokens.RemoveTokens(c, true, true)

	logoutResult = ProcessUiResult(c, logoutResult, errResp.WithUI("/auth/login", "/auth/login", "Login"), span, log, h.LogoutWEBCounter, "logout_web")
	if logoutResult == nil {
		return
	}

	// Increment the counter for web logout
	h.LogoutWEBCounter.WithLabelValues("success").Inc()
	c.Redirect(http.StatusSeeOther, "/auth/login") // Redirect to login page after logout
}

// HandleLogoutAPI godoc
// @Summary      Logout user API
// @Description  Logs out the user by revoking their session using session token or refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        logoutData  body      auth.LogoutData  true  "Logout data containing session token or refresh token"
// @Success      200  {object}  auth.LogoutResult "Logout result"
// @Failure      400  {object}  models.ErrorResponse "Invalid request"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /api/auth/s/logout [post]
func (h *LogoutHandler) HandleLogoutAPI(c *gin.Context) {
	h.LogoutAPICounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "logout_api")
	defer span.End()

	var logoutData auth.LogoutData
	if err := c.ShouldBindJSON(&logoutData); err != nil {
		h.LogoutAPICounter.WithLabelValues("invalid_request").Inc()
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.GenericErrorMessage, "Invalid request", http.StatusBadRequest, err))
		return
	}

	logoutData.SessionToken = strings.TrimSpace(logoutData.SessionToken)
	logoutData.RefreshToken = strings.TrimSpace(logoutData.RefreshToken)

	if logoutData.SessionToken == "" && logoutData.RefreshToken == "" {
		h.LogoutAPICounter.WithLabelValues("invalid_request").Inc()
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.GenericErrorMessage, "Either session token or refresh token must be provided", http.StatusBadRequest, nil))
		return
	}

	logoutResult, errResp := auth.HandleLogout(ctx, log, logoutData)

	logoutResult = ProcessAPIResult(c, logoutResult, errResp, span, log, h.LogoutAPICounter, "logout_api")
	if logoutResult == nil {
		return
	}

	// Increment the counter for API logout
	h.LogoutAPICounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, logoutResult)
}
