package auth

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"go.uber.org/zap"
)

// LogoutData represents the data required for logging out a user.
//
// It includes the access token and refresh token.
// Use strings.TrimSpace to ensure no leading or trailing spaces are present.
type LogoutData struct {
	SessionToken string `json:"session_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HandleLogout processes the logout request by revoking the session
// using the provided session token or refresh token.
// It returns a LogoutResult indicating success or failure, and an ErrorResponse if an error occurs.
//
// NOTE: One of SessionToken or RefreshToken must be provided.
func HandleLogout(ctx context.Context, log *zap.Logger, logoutData LogoutData) (*LogoutResult, *models.ErrorResponse) {
	log.Debug("Handling logout")

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error("Failed to begin transaction", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to begin transaction", http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)
	log.Debug("Transaction begun successfully")

	if logoutData.SessionToken != "" {
		log.Debug("Revoking using session token")
		tokenHash, _ := tokens.HashTokens(&tokens.Tokens{
			SessionToken: logoutData.SessionToken,
		})
		err = q.DeleteSessionByToken(ctx, tokenHash)
	} else {
		// This condition is commented because we assume that the session token OR the refresh token is always provided.
		// } else if logoutData.RefreshToken != "" {
		log.Debug("Revoking using refresh token")
		_, tokenHash := tokens.HashTokens(&tokens.Tokens{
			RefreshToken: logoutData.RefreshToken,
		})
		err = q.DeleteSessionByRefreshToken(ctx, tokenHash)
	}

	if err != nil {
		log.Error("Failed to revoke session", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to revoke session", http.StatusInternalServerError, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction", http.StatusInternalServerError, err)
	}

	log.Debug("Session revoked successfully")
	return &LogoutResult{
		Success: true,
		Message: "Session revoked successfully",
	}, nil
}
