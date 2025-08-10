package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nbrglm/auth-platform/db"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/notifications"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"go.uber.org/zap"
)

type SendVerificationEmailData struct {
	Email string `json:"email" binding:"required,email"`
}

type SendVerificationEmailResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SendVerificationEmail(ctx context.Context, log *zap.Logger, input SendVerificationEmailData) (*SendVerificationEmailResult, *models.ErrorResponse) {
	log.Debug("Sending verification email", zap.String("email", input.Email))

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to begin transaction!", http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)

	// Check if the user exists
	user, err := q.GetUserByEmail(ctx, input.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.NewErrorResponse("It seems the user with the provided email does not exist! Please check the email and try again.", "No user exists with the provided email address.", http.StatusNotFound, nil)
	}
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to retrieve user information!", http.StatusInternalServerError, err)
	}

	if user.EmailVerified {
		return nil, models.NewErrorResponse("The email is already verified! No need to send a verification email again.", "Email already verified!", http.StatusBadRequest, nil)
	}

	// Generate a verification token
	token, hash, err := tokens.GenerateEmailVerificationToken()
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to generate verification token!", http.StatusInternalServerError, err)
	}

	// Generate the token ID
	tokenId, err := uuid.NewV7()
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to generate token ID!", http.StatusInternalServerError, err)
	}

	// Insert the verification token into the database
	newToken, err := q.NewVerificationToken(ctx, db.NewVerificationTokenParams{
		ID:        tokenId,
		UserID:    user.ID,
		Type:      string(tokens.EmailVerificationToken),
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		},
	})

	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to store verification token!", http.StatusInternalServerError, err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to commit transaction!", http.StatusInternalServerError, err)
	}

	// Send the verification email
	if err := notifications.SendWelcomeEmail(ctx, notifications.SendWelcomeEmailParams{
		User: struct {
			Email     string
			FirstName *string
			LastName  *string
		}{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		VerificationToken: token,
		ExpiresAt:         newToken.ExpiresAt.Time,
	}); err != nil {
		return nil, models.NewErrorResponse("An error occurred while sending the verification email. Please try again later.", "Failed to send verification email!", http.StatusInternalServerError, err)
	}

	return &SendVerificationEmailResult{
		Success: true,
		Message: "Verification email sent successfully. Please check your inbox.",
	}, nil
}

// VerifyEmailTokenData defines the input structure for verifying an email token.
type VerifyEmailTokenData struct {
	Token string `json:"token" form:"token" binding:"required"`
}

// VerifyEmailTokenResult defines the response structure for verifying an email token.
type VerifyEmailTokenResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func VerifyEmailToken(ctx context.Context, log *zap.Logger, input VerifyEmailTokenData) (*VerifyEmailTokenResult, *models.ErrorResponse) {
	log.Debug("Verifying email token", zap.String("token", input.Token))

	// Begin TX
	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to begin transaction!", http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)

	// Hash the token first
	hash := tokens.HashEmailVerificationToken(input.Token)

	// Fetch the verification token from the database
	token, err := q.GetVerificationTokenByHash(ctx, hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.NewErrorResponse("Invalid token! Please check the token and try again.", "No verification token found with the provided token hash.", http.StatusNotFound, nil)
	}
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to retrieve verification token!", http.StatusInternalServerError, err)
	}
	if token.Type != string(tokens.EmailVerificationToken) {
		return nil, models.NewErrorResponse("Invalid token! Please check the token and try again.", "The provided token is not a valid email verification token.", http.StatusBadRequest, nil)
	}
	if !token.ExpiresAt.Valid || token.ExpiresAt.Time.Before(time.Now()) {
		return nil, models.NewErrorResponse("The token has expired! Please request a new verification email.", "The provided token has expired.", http.StatusBadRequest, nil)
	}

	// Since a token was found and is valid, we can proceed to mark the user's email as verified
	err = q.MarkUserEmailVerified(ctx, token.UserID)
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to mark user email as verified!", http.StatusInternalServerError, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to commit transaction!", http.StatusInternalServerError, err)
	}
	log.Debug("Email verified successfully", zap.String("userID", token.UserID.String()))
	return &VerifyEmailTokenResult{
		Success: true,
		Message: "Email verified successfully!",
	}, nil
}
