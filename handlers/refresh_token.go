package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nbrglm/auth-platform/db"
	"github.com/nbrglm/auth-platform/internal"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type RefreshTokenHandler struct {
	// Track the number of refresh token requests made by the UI.
	RefreshWEBCounter *prometheus.CounterVec

	// Track the number of refresh token requests made by the API.
	RefreshAPICounter *prometheus.CounterVec
}

func NewRefreshTokenHandler() *RefreshTokenHandler {
	return &RefreshTokenHandler{
		RefreshWEBCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_refresh_web_requests",
				Help:      "Total number of user refresh token requests from the web",
			},
			[]string{"status"},
		),
		RefreshAPICounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_refresh_api_requests",
				Help:      "Total number of user refresh token requests from the API",
			},
			[]string{"status"},
		),
	}
}

func (h *RefreshTokenHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.RefreshAPICounter)
	metrics.Collectors = append(metrics.Collectors, h.RefreshWEBCounter)

	engine.POST("/api/auth/s/refresh", middlewares.RateLimitAPIMiddleware(), h.HandleRefreshTokenAPI)
	engine.GET("/auth/s/refresh", middlewares.RateLimitUIAuthenticatedMiddleware(), h.HandleRefreshTokenWEB)
}

type RefreshTokenData struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// HandleRefreshTokenAPI godoc
// @Summary Refresh Token API
// @Description Handles token refresh requests via API.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body auth.RefreshTokenData true "Refresh Token Data"
// @Success 200 {object} auth.SessionRefreshResult "Refresh Token Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/s/refresh [post]
func (h *RefreshTokenHandler) HandleRefreshTokenAPI(c *gin.Context) {
	h.RefreshAPICounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "refresh_token_api")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var refreshData RefreshTokenData
	if err := c.ShouldBindJSON(&refreshData); err != nil {
		h.RefreshAPICounter.WithLabelValues("invalid_input").Inc()
		log.Debug("Invalid input data", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request data", "Please check your input and try again.", http.StatusBadRequest, nil).Filter())
		return
	}

	result, err := handleRefresh(ctx, log, refreshData)
	result = ProcessAPIResult(c, result, err, span, log, h.RefreshAPICounter, "refresh_token_api")
	if result == nil {
		return
	}
	if result.ShouldLogin {
		h.RefreshAPICounter.WithLabelValues("invalid_token").Inc()
		// Only place where we break convention,
		// since we cannot return a 401 Unauthorized here as the requester is an application's backend,
		// we return a 200 with a flag indicating the user should login.
		c.JSON(http.StatusOK, result)
		return
	}

	h.RefreshAPICounter.WithLabelValues("success").Inc()
	span.SetStatus(codes.Ok, "Tokens refreshed successfully")
	c.JSON(http.StatusOK, result.Tokens)
}

func (h *RefreshTokenHandler) HandleRefreshTokenWEB(c *gin.Context) {
	h.RefreshWEBCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "refresh_token_web")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	log.Debug("Checking refresh status for user",
		zap.Bool("refresh_needed", c.GetBool(middlewares.CtxSessionRefreshKey)),
		zap.Bool("session_exists", c.GetBool(middlewares.CtxSessionExistsKey)))

	if c.GetBool(middlewares.CtxSessionExistsKey) {
		h.RefreshWEBCounter.WithLabelValues("no_refresh_needed").Inc()
		u := middlewares.GetRedirectURLOriginalOrFallback(c)
		c.Redirect(http.StatusSeeOther, u)
		return
	} else if !c.GetBool(middlewares.CtxSessionRefreshKey) {
		h.RefreshWEBCounter.WithLabelValues("invalid_token").Inc()
		u := middlewares.GetRedirectURLWithReturnTo(c, "/auth/login")
		c.Redirect(http.StatusSeeOther, u)
		return
	}

	refreshData := RefreshTokenData{
		RefreshToken: c.GetString(middlewares.CtxSessionRefreshTokenKey),
	}

	result, err := handleRefresh(ctx, log, refreshData)
	log.Debug("Refresh tokens result", zap.Any("result", result), zap.Error(err))

	if err != nil {
		// To prevent an infinite loop, we REMOVE the refreshToken from the cookies
		tokens.RemoveTokens(c, false, true)
		// The actual processing of the error is done in ProcessUiResult,
		// which will also set the appropriate HTTP status code and response body.
	}

	// The WithUI's params other than the returnURL do not matter as we don't show an error...
	result = ProcessUiResult(c, result, err.WithUI("", middlewares.GetRedirectURLWithReturnTo(c, "/auth/login"), ""), span, log, h.RefreshWEBCounter, "refresh_token_web")
	if result == nil {
		return
	}

	if result.ShouldLogin {
		// To prevent an infinite loop, we REMOVE the refreshToken from the cookies
		tokens.RemoveTokens(c, false, true)

		h.RefreshWEBCounter.WithLabelValues("invalid_token").Inc()
		c.Redirect(http.StatusSeeOther, middlewares.GetRedirectURLWithReturnTo(c, "/auth/login"))
		return
	}

	tokens.SetTokens(c, result.Tokens)
	h.RefreshWEBCounter.WithLabelValues("success").Inc()
	span.SetStatus(codes.Ok, "Tokens refreshed successfully")
	c.Redirect(http.StatusSeeOther, middlewares.GetRedirectURLWithReturnTo(c, "/auth/login"))
}

type SessionRefreshResult struct {
	Tokens      *tokens.Tokens
	ShouldLogin bool `json:"shouldLogin"`
}

func handleRefresh(ctx context.Context, log *zap.Logger, refreshData RefreshTokenData) (*SessionRefreshResult, *models.ErrorResponse) {
	log.Debug("Starting token refresh process", zap.String("refreshToken", refreshData.RefreshToken))

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Debug("Failed to begin transaction", zap.Error(err))
		return nil, models.NewErrorResponse(
			models.GenericErrorMessage,
			"Failed to begin transaction",
			http.StatusInternalServerError,
			err,
		)
	}

	q := store.Querier.WithTx(tx)
	log.Debug("Transaction begun successfully")

	_, refreshTokenHash := tokens.HashTokens(&tokens.Tokens{
		RefreshToken: refreshData.RefreshToken,
	})
	log.Debug("Refresh token hashed")

	session, err := q.GetSessionByRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("No session found for refresh token")
			return &SessionRefreshResult{
				ShouldLogin: true,
			}, nil
		}
		log.Debug("Error retrieving session", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Unable to retrieve session", http.StatusInternalServerError, err)
	}
	log.Debug("Session retrieved successfully", zap.String("sessionID", session.ID.String()))

	revoked, err := tokens.HasTokenBeenRevoked(ctx, q, session.ID)
	if err != nil {
		log.Debug("Error checking token revocation status", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Unable to check token status", http.StatusInternalServerError, err)
	}

	// If the token has been revoked, return an error
	if revoked {
		log.Debug("Token has been revoked", zap.String("sessionID", session.ID.String()))
		return &SessionRefreshResult{
			ShouldLogin: true,
		}, nil
	}
	log.Debug("Token is valid and not revoked")

	newTokenInfo, err := q.GetInfoForSessionRefresh(ctx, db.GetInfoForSessionRefreshParams{
		UserID: session.UserID,
		OrgID:  session.OrgID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("User or org info not found for session",
				zap.String("userID", session.UserID.String()),
				zap.String("orgID", session.OrgID.String()))
			// Either user or org associated with this session was not found, or the user has been banned.
			return &SessionRefreshResult{
				ShouldLogin: true,
			}, nil
		}
		log.Debug("Error retrieving session info", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Unable to retrieve session info", http.StatusInternalServerError, err)
	}
	log.Debug("User and org info retrieved successfully")

	avatarUrl := ""
	if newTokenInfo.UserAvatarUrl != nil {
		avatarUrl = *newTokenInfo.UserAvatarUrl
	}

	log.Debug("Generating new token pair")
	newTokenPair, err := tokens.RefreshSessionTokens(session, tokens.AuthPlatformClaims{
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
		log.Debug("Error generating new tokens", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Unable to generate new tokens", http.StatusInternalServerError, err)
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
		log.Debug("Error refreshing session in database", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Unable to refresh session", http.StatusInternalServerError, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction!", http.StatusInternalServerError, err)
	}

	log.Debug("Session refreshed successfully", zap.String("sessionID", session.ID.String()))
	return &SessionRefreshResult{
		Tokens: newTokenPair,
	}, nil
}
