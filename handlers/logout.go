package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/nbrglm/nexeres/internal"
	"github.com/nbrglm/nexeres/internal/metrics"
	"github.com/nbrglm/nexeres/internal/middlewares"
	"github.com/nbrglm/nexeres/internal/models"
	"github.com/nbrglm/nexeres/internal/store"
	"github.com/nbrglm/nexeres/internal/tokens"
	"github.com/nbrglm/nexeres/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type LogoutHandler struct {
	LogoutCounter *prometheus.CounterVec
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{
		LogoutCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "auth",
				Name:      "user_logout_requests",
				Help:      "Total number of user logout requests",
			},
			[]string{"status"},
		),
	}
}

func (h *LogoutHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.LogoutCounter)
	engine.POST("/api/auth/logout", middlewares.RequireAuth(middlewares.AuthModeEither), h.HandleLogout)
}

type LogoutResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HandleLogout godoc
// @Summary      Logout user
// @Description  Logs out the user by revoking their session using session token or refresh token. Requires atleast one of the tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param X-NEXERES-Session-Token header string false "Session token"
// @Param X-NEXERES-Refresh-Token header string false "Refresh token"
// @Success      200  {object}  LogoutResult "Logout result"
// @Failure      400  {object}  models.ErrorResponse "Bad Request - Invalid or missing tokens"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized - Invalid or expired tokens"
// @Failure      500  {object}  models.ErrorResponse "Internal server error
// @Router       /api/auth/logout [post]
func (h *LogoutHandler) HandleLogout(c *gin.Context) {
	h.LogoutCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "logout")
	defer span.End()

	sessionToken := c.GetString(middlewares.CtxSessionToken)
	refreshToken := c.GetString(middlewares.CtxRefreshToken)

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to begin transaction", http.StatusInternalServerError, err), span, log, h.LogoutCounter, "logout")
		return
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)
	log.Debug("Transaction begun successfully")

	if sessionToken != "" {
		log.Debug("Handling logout request with session token")
		tokenHash, _ := tokens.HashTokens(&tokens.Tokens{
			SessionToken: sessionToken,
		})
		err = q.DeleteSessionByToken(ctx, tokenHash)
	} else {
		log.Debug("Handling logout request with refresh token")
		_, tokenHash := tokens.HashTokens(&tokens.Tokens{
			RefreshToken: refreshToken,
		})
		err = q.DeleteSessionByRefreshToken(ctx, tokenHash)
	}

	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to revoke session", http.StatusInternalServerError, err), span, log, h.LogoutCounter, "logout")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction", http.StatusInternalServerError, err), span, log, h.LogoutCounter, "logout")
		return
	}

	log.Debug("Session revoked successfully")
	c.JSON(http.StatusOK, &LogoutResult{
		Success: true,
		Message: "Session revoked successfully",
	})
}
