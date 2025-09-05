package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nbrglm/nexeres/db"
	"github.com/nbrglm/nexeres/internal"
	"github.com/nbrglm/nexeres/internal/metrics"
	"github.com/nbrglm/nexeres/internal/models"
	"github.com/nbrglm/nexeres/internal/notifications"
	"github.com/nbrglm/nexeres/internal/store"
	"github.com/nbrglm/nexeres/internal/tokens"
	"github.com/nbrglm/nexeres/utils"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type VerifyEmailHandler struct {
	SendEmailCounter   *prometheus.CounterVec
	VerifyEmailCounter *prometheus.CounterVec
}

func NewVerifyEmailHandler() *VerifyEmailHandler {
	return &VerifyEmailHandler{
		SendEmailCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "auth",
				Name:      "verify_email_send_requests",
				Help:      "Total number of requests to send verification email",
			},
			[]string{"status"},
		),
		VerifyEmailCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "auth",
				Name:      "verify_email_requests",
				Help:      "Total number of requests to verify email with token",
			},
			[]string{"status"},
		),
	}
}

func (h *VerifyEmailHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.SendEmailCounter)
	metrics.Collectors = append(metrics.Collectors, h.VerifyEmailCounter)

	engine.POST("/api/auth/verify-email/send", h.HandleSendVerificationEmail)
	engine.POST("/api/auth/verify-email/verify", h.HandleVerifyEmailToken)
}

type SendVerificationEmailData struct {
	Email string `json:"email" binding:"required,email"`
}

type SendVerificationEmailResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HandleSendVerificationEmail godoc
// @Summary Send Verification Email
// @Description Sends a verification email to the user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body SendVerificationEmailData true "Send Verification Email Data"
// @Success 200 {object} SendVerificationEmailResult "Send Verification Email Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid Input or User does not exist"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/verify-email/send [post]
func (h *VerifyEmailHandler) HandleSendVerificationEmail(c *gin.Context) {
	h.SendEmailCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "send_verification_email_api")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var input SendVerificationEmailData
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Invalid request data. Please check your input and try again.", "Failed to bind JSON!", http.StatusBadRequest, nil), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to begin transaction!", http.StatusInternalServerError, err), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)

	// Check if the user exists
	user, err := q.GetUserByEmail(ctx, input.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		utils.ProcessError(c, models.NewErrorResponse("It seems the user with the provided email does not exist! Please check the email and try again.", "No user exists with the provided email address.", http.StatusBadRequest, nil), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to retrieve user information!", http.StatusInternalServerError, err), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}

	if user.EmailVerified {
		utils.ProcessError(c, models.NewErrorResponse("The email is already verified!", "Email already verified!", http.StatusBadRequest, nil), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}

	// Generate a verification token
	token, hash, err := tokens.GenerateEmailVerificationToken()
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to generate verification token!", http.StatusInternalServerError, err), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}

	// Generate the token ID
	tokenId, err := uuid.NewV7()
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to generate token ID!", http.StatusInternalServerError, err), span, log, h.SendEmailCounter, "send_verification_email")
		return
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
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to insert verification token into the database!", http.StatusInternalServerError, err), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction!", http.StatusInternalServerError, err), span, log, h.SendEmailCounter, "send_verification_email")
		return
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
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to send verification email!", http.StatusInternalServerError, err), span, log, h.SendEmailCounter, "send_verification_email")
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, SendVerificationEmailResult{
		Success: true,
		Message: "Verification email sent successfully!",
	})
}

type VerifyEmailTokenData struct {
	Token string `json:"token" binding:"required"`
}

type VerifyEmailTokenResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HandleVerifyEmailToken godoc
// @Summary Verify Email Token
// @Description Verifies the email using the provided token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body VerifyEmailTokenData true "Verify Email Token Data"
// @Success 200 {object} VerifyEmailTokenResult "Verify Email Token Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request - Invalid Input or Token"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/verify-email/verify [post]
func (h *VerifyEmailHandler) HandleVerifyEmailToken(c *gin.Context) {
	h.VerifyEmailCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "verify_email_token")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var input VerifyEmailTokenData
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Invalid request data. Please check your input and try again.", "Failed to bind JSON!", http.StatusBadRequest, nil), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to begin transaction!", http.StatusInternalServerError, err), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)

	// Hash the token first
	hash := tokens.HashEmailVerificationToken(input.Token)

	// Fetch the verification token from the database
	token, err := q.GetVerificationTokenByHash(ctx, hash)
	if errors.Is(err, pgx.ErrNoRows) {
		utils.ProcessError(c, models.NewErrorResponse("Invalid token! Please check the token and try again.", "No verification token found with the provided token hash.", http.StatusBadRequest, nil), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to retrieve verification token!", http.StatusInternalServerError, err), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}
	if token.Type != string(tokens.EmailVerificationToken) {
		utils.ProcessError(c, models.NewErrorResponse("Invalid token! Please check the token and try again.", "The provided token is not a valid email verification token.", http.StatusBadRequest, nil), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}
	if !token.ExpiresAt.Valid || token.ExpiresAt.Time.Before(time.Now()) {
		utils.ProcessError(c, models.NewErrorResponse("The token has expired! Please request a new verification email.", "The provided token has expired.", http.StatusBadRequest, nil), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}

	// Mark the user's email as verified, since a valid token was found
	err = q.MarkUserEmailVerified(ctx, token.UserID)
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to mark user email as verified!", http.StatusInternalServerError, err), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}
	if err := tx.Commit(ctx); err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction!", http.StatusInternalServerError, err), span, log, h.VerifyEmailCounter, "verify_email_token")
		return
	}
	log.Debug("Email verified successfully", zap.String("userID", token.UserID.String()))

	c.JSON(http.StatusOK, VerifyEmailTokenResult{
		Success: true,
		Message: "Email verified successfully!",
	})
}
