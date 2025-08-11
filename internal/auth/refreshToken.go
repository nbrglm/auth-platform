package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nbrglm/auth-platform/db"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"go.uber.org/zap"
)

type RefreshTokenData struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SessionRefreshResult struct {
	Tokens      *tokens.Tokens
	ShouldLogin bool `json:"shouldLogin"`
}

// HandleRefresh processes the refresh token request.
// It checks the validity of the refresh token, retrieves the session, and generates new tokens if valid.
func HandleRefresh(ctx context.Context, log *zap.Logger, refreshData RefreshTokenData) (*SessionRefreshResult, *models.ErrorResponse) {
	log.Debug("Starting token refresh process")

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
	defer tx.Rollback(ctx)

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
