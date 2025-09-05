package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nbrglm/nexeres/db"
	"github.com/nbrglm/nexeres/internal"
	"github.com/nbrglm/nexeres/internal/metrics"
	"github.com/nbrglm/nexeres/internal/middlewares"
	"github.com/nbrglm/nexeres/internal/models"
	"github.com/nbrglm/nexeres/internal/store"
	"github.com/nbrglm/nexeres/internal/tokens"
	"github.com/nbrglm/nexeres/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type RefreshTokenHandler struct {
	RefreshTokenCounter *prometheus.CounterVec
}

func NewRefreshTokenHandler() *RefreshTokenHandler {
	return &RefreshTokenHandler{
		RefreshTokenCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "auth",
				Name:      "user_refresh_token_requests",
				Help:      "Total number of user refresh token requests",
			},
			[]string{"status"},
		),
	}
}

func (h *RefreshTokenHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.RefreshTokenCounter)
	engine.POST("/api/auth/refresh", middlewares.RequireAuth(middlewares.AuthModeRefresh), h.HandleRefreshToken)
}

type RefreshTokenResult struct {
	Tokens *tokens.Tokens `json:"tokens"`
}

// HandleRefreshToken godoc
// @Summary Refresh Token
// @Description Handles token refresh requests.
// @Tags Auth
// @Accept json
// @Produce json
// @Param X-NEXERES-Refresh-Token header string true "Refresh token"
// @Success 200 {object} RefreshTokenResult "New tokens"
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid or missing tokens"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Invalid or expired tokens - Proceed to Login"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/auth/refresh [post]
func (h *RefreshTokenHandler) HandleRefreshToken(c *gin.Context) {
	h.RefreshTokenCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "refresh_token_api")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	refreshToken := c.GetString(middlewares.CtxRefreshToken)
	// we don't check if refreshToken is empty because the RequireAuth middleware ensures it's present

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to begin transaction", http.StatusInternalServerError, err), span, log, h.RefreshTokenCounter, "refresh_token")
		return
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)
	log.Debug("Transaction begun successfully")

	_, refreshTokenHash := tokens.HashTokens(&tokens.Tokens{
		RefreshToken: refreshToken,
	})
	log.Debug("Refresh token hashed")

	session, err := q.GetSessionByRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.ProcessError(c, models.NewErrorResponse("Invalid refresh token! Please login again.", "No session found for refresh token", http.StatusUnauthorized, nil), span, log, h.RefreshTokenCounter, "refresh_token")
			return
		}

		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Unable to retrieve session", http.StatusInternalServerError, err), span, log, h.RefreshTokenCounter, "refresh_token")
		return
	}
	log.Debug("Session retrieved successfully", zap.String("sessionID", session.ID.String()))

	// We DO NOT CHECK if the token has been revoked here as:
	// The revocation is done by deleting the session from the database.
	// So if the session exists, the token is valid and not revoked.
	// Thus, if we are here, it means the session exists and the token is valid.
	// Using the method `tokens.HasTokenBeenRevoked`, which checks if no session exists for the given session ID,
	// would be redundant.

	newTokenInfo, err := q.GetInfoForSessionRefresh(ctx, db.GetInfoForSessionRefreshParams{
		UserID: session.UserID,
		OrgID:  session.OrgID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("No user or organization found for session", zap.String("sessionID", session.ID.String()))
			utils.ProcessError(c, models.NewErrorResponse("User or organization not found! Please login again.", "No user or organization found for session", http.StatusUnauthorized, nil), span, log, h.RefreshTokenCounter, "refresh_token")
			return
		}

		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Unable to retrieve user or organization info", http.StatusInternalServerError, err), span, log, h.RefreshTokenCounter, "refresh_token")
		return
	}
	log.Debug("User and organization info retrieved successfully", zap.String("userID", session.UserID.String()), zap.String("orgSlug", newTokenInfo.OrgSlug))

	avatarUrl := ""
	if newTokenInfo.UserAvatarUrl != nil {
		avatarUrl = *newTokenInfo.UserAvatarUrl
	}

	log.Debug("Generating new tokens for user", zap.String("orgSlug", newTokenInfo.OrgSlug))

	newTokenPair, err := tokens.RefreshSessionTokens(session, tokens.NexeresClaims{
		OrgSlug: newTokenInfo.OrgSlug,
		OrgName: newTokenInfo.OrgName,
		OrgId:   session.OrgID.String(),

		Email:         newTokenInfo.UserEmail,
		EmailVerified: newTokenInfo.UserEmailVerified,
		UserFname:     *newTokenInfo.UserFname,
		UserLname:     *newTokenInfo.UserLname,
		UserAvatarURL: avatarUrl,
		UserOrgRole:   newTokenInfo.UserOrgRole,
	})
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Unable to generate new tokens", http.StatusInternalServerError, err), span, log, h.RefreshTokenCounter, "refresh_token")
		return
	}
	log.Debug("New token pair generated successfully")

	newSessionTokenHash, newRefreshTokenHash := tokens.HashTokens(newTokenPair)

	_, err = q.RefreshSession(ctx, db.RefreshSessionParams{
		ID:               session.ID,
		RefreshTokenHash: &newRefreshTokenHash,
		TokenHash:        &newSessionTokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  newTokenPair.RefreshTokenExpiry,
			Valid: true,
		},
	})

	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Unable to refresh session", http.StatusInternalServerError, err), span, log, h.RefreshTokenCounter, "refresh_token")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction!", http.StatusInternalServerError, err), span, log, h.RefreshTokenCounter, "refresh_token")
		return
	}

	h.RefreshTokenCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, RefreshTokenResult{
		Tokens: newTokenPair,
	})
}
