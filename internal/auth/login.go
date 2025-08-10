package auth

import (
	"context"
	"errors"
	"net/http"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/db"
	"github.com/nbrglm/auth-platform/internal/cache"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/password"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"github.com/nbrglm/auth-platform/opts"
	"github.com/nbrglm/auth-platform/utils"
	"go.uber.org/zap"
)

type UserLoginData struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=32"`
	// Not to include in the input data from client

	// The IP address of the user making the request, used for logging and security purposes
	IPAddress netip.Addr `json:"-"`

	// The user agent of the client making the request, used for logging and security purposes
	UserAgent string `json:"-"`

	// Only used when the UI sets a flow ID, which contains the returnTo URL
	// This is used to continue the flow after the user has logged in.
	FlowID *string `json:"flowId" form:"flowId" binding:"omitempty"`
}

type UserLoginResult struct {
	Message string `json:"message"`
	tokens.Tokens
	SentVerificationEmail bool    `json:"sentVerificationEmail"`
	FlowID                *string `json:"flowId"` // For multitenant flows, the flow ID to continue the flow
}

func HandleLogin(ctx context.Context, log *zap.Logger, input UserLoginData) (*UserLoginResult, *models.ErrorResponse) {
	log.Debug("Login attempt started", zap.String("email", input.Email))

	if config.Multitenancy.Enable {
		log.Debug("Multitenancy enabled, handling multitenant login")
		return handleMultitenantLogin(ctx, log, input)
	}

	log.Debug("Processing single-tenant login")
	// For single-tenant login, we can directly validate the user's credentials
	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error("Failed to begin transaction", zap.Error(err))
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to begin transaction!", http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)

	log.Debug("Retrieving user login information")
	user, err := q.GetLoginInfoForUser(ctx, input.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Debug("User not found", zap.String("email", input.Email))
		return nil, models.NewErrorResponse("Invalid email or password! Please try again.", "User not found!", http.StatusUnauthorized, nil)
	}
	if err != nil {
		log.Error("Failed to retrieve user information", zap.Error(err))
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to retrieve user information!", http.StatusInternalServerError, err)
	}

	if !user.EmailVerified {
		log.Debug("Email not verified, sending verification email", zap.String("email", user.Email))
		// Send a verification email if the email is not verified
		result, err := SendVerificationEmail(ctx, log, SendVerificationEmailData{
			Email: user.Email,
		})

		if err != nil {
			log.Error("Failed to send verification email", zap.Error(err))
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to send verification email!", http.StatusInternalServerError, err)
		}
		if !result.Success {
			log.Debug("Verification email sending failed")
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to send verification email!", http.StatusInternalServerError, nil)
		}
		log.Debug("Verification email sent successfully")
		return &UserLoginResult{
			Message:               "A verification email has been sent to your email address. Please verify your email before logging in.",
			SentVerificationEmail: true,
		}, nil
	}

	log.Debug("Verifying password")
	if !password.VerifyPasswordMatch(*user.PasswordHash, input.Password) {
		log.Debug("Password mismatch", zap.String("email", input.Email))
		return nil, models.NewErrorResponse("Invalid email or password! Please try again.", "Password mismatch!", http.StatusUnauthorized, nil)
	}

	avatarUrl := ""
	if user.AvatarUrl != nil {
		avatarUrl = *user.AvatarUrl
	}

	log.Debug("Generating tokens")
	result, err := tokens.GenerateTokens(user.ID, tokens.AuthPlatformClaims{
		OrgSlug: opts.DefaultOrgSlug,
		OrgName: opts.DefaultOrgName,
		OrgId:   opts.DefaultOrgId,

		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		UserFname:     *user.FirstName,
		UserLname:     *user.LastName,
		UserAvatarURL: avatarUrl,
		UserOrgRole:   "member",
	})
	if err != nil {
		log.Error("Failed to generate tokens", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to generate token pair!", http.StatusInternalServerError, err)
	}

	newSessionTokenHash, newRefreshTokenHash := tokens.HashTokens(result)
	log.Debug("Creating session in database", zap.String("sessionId", result.SessionId.String()))

	_, err = q.CreateSession(ctx, db.CreateSessionParams{
		ID:               result.SessionId,
		UserID:           user.ID,
		OrgID:            uuid.MustParse(opts.DefaultOrgId),
		TokenHash:        newSessionTokenHash,
		RefreshTokenHash: newRefreshTokenHash,
		MfaVerified:      false,
		MfaVerifiedAt: pgtype.Timestamptz{
			Valid: false,
		},
		IpAddress: input.IPAddress,
		UserAgent: input.UserAgent,
		ExpiresAt: pgtype.Timestamptz{
			Time:  result.RefreshTokenExpiry,
			Valid: true,
		},
	})

	if err != nil {
		log.Error("Failed to create session in the db", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to create session in the db!", http.StatusInternalServerError, err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction!", http.StatusInternalServerError, err)
	}

	log.Debug("Login successful", zap.String("email", user.Email), zap.String("sessionId", result.SessionId.String()))
	return &UserLoginResult{
		Message:               "Login successful",
		Tokens:                *result,
		SentVerificationEmail: false,
	}, nil
}

// handleMultitenantLogin handles the login flow for multitenant setups.
//
// The flow is as follows:
// 1. Validate the email and password
// 2. If valid, fetch the organizations the user belong to.
// 3. If the user belongs to a single organization, create a session for that organization and return the session and refresh tokens.
// 4. If the user belongs to multiple organizations, create a new flow to select the organization, store it, and redirect the user to /auth/select-org?flow=<flow-id>
func handleMultitenantLogin(ctx context.Context, log *zap.Logger, data UserLoginData) (*UserLoginResult, *models.ErrorResponse) {
	_, err := utils.GetDomainFromEmail(data.Email)
	if err != nil {
		return nil, models.NewErrorResponse("Invalid request! Please input a valid email and try again.", "Invalid email domain!", http.StatusBadRequest, nil)
	}

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to begin transaction!", http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	q := store.Querier.WithTx(tx)

	user, err := q.GetLoginInfoForUser(ctx, data.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.NewErrorResponse("Invalid email or password! Please try again.", "User not found!", http.StatusUnauthorized, nil)
	}
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to retrieve user information!", http.StatusInternalServerError, err)
	}

	if !user.EmailVerified {
		// Send a verification email if the email is not verified
		result, err := SendVerificationEmail(ctx, log, SendVerificationEmailData{
			Email: user.Email,
		})

		if err != nil {
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to send verification email!", http.StatusInternalServerError, err)
		}
		if !result.Success {
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to send verification email!", http.StatusInternalServerError, nil)
		}
		return &UserLoginResult{
			Message:               "A verification email has been sent to your email address. Please verify your email before logging in.",
			SentVerificationEmail: true,
		}, nil
	}

	if !password.VerifyPasswordMatch(*user.PasswordHash, data.Password) {
		return nil, models.NewErrorResponse("Invalid email or password! Please try again.", "Password mismatch!", http.StatusUnauthorized, nil)
	}
	// Fetch the organizations the user belongs to
	orgs, err := q.GetUserOrgsByEmail(ctx, &user.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.NewErrorResponse("You do not belong to any organization! Please contact your administrator.", "No organizations found for the user!", http.StatusNotFound, nil)
	}
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to retrieve user organizations!", http.StatusInternalServerError, err)
	}
	if len(orgs) == 0 {
		return nil, models.NewErrorResponse("You do not belong to any organization! Please contact your administrator.", "No organizations found for the user!", http.StatusNotFound, nil)
	}
	if len(orgs) == 1 {
		avatarUrl := ""
		if user.AvatarUrl != nil {
			avatarUrl = *user.AvatarUrl
		}
		// If the user belongs to a single organization, create a session for that organization
		result, err := tokens.GenerateTokens(user.ID, tokens.AuthPlatformClaims{
			OrgSlug: orgs[0].Org.Slug,
			OrgName: orgs[0].Org.Name,
			OrgId:   orgs[0].Org.ID.String(),

			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			UserFname:     *user.FirstName,
			UserLname:     *user.LastName,
			UserAvatarURL: avatarUrl,
			UserOrgRole:   orgs[0].UserOrg.Role,
		})
		if err != nil {
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to generate token pair!", http.StatusInternalServerError, err)
		}

		newSessionTokenHash, newRefreshTokenHash := tokens.HashTokens(result)

		_, err = q.CreateSession(ctx, db.CreateSessionParams{
			ID:               result.SessionId,
			UserID:           user.ID,
			OrgID:            orgs[0].Org.ID,
			TokenHash:        newSessionTokenHash,
			RefreshTokenHash: newRefreshTokenHash,
			MfaVerified:      false,
			MfaVerifiedAt: pgtype.Timestamptz{
				Valid: false,
			},
			IpAddress: data.IPAddress,
			UserAgent: data.UserAgent,
			ExpiresAt: pgtype.Timestamptz{
				Time:  result.RefreshTokenExpiry,
				Valid: true,
			},
		})

		if err != nil {
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to create session in the db!", http.StatusInternalServerError, err)
		}

		if err := tx.Commit(ctx); err != nil {
			log.Error("Failed to commit transaction", zap.Error(err))
			return nil, models.NewErrorResponse(models.GenericErrorMessage, "Failed to commit transaction!", http.StatusInternalServerError, err)
		}

		return &UserLoginResult{
			Message:               "Login successful",
			Tokens:                *result,
			SentVerificationEmail: false,
		}, nil
	}

	// If the user belongs to multiple organizations, create a new flow to select the organization
	organizations := make([]db.Org, len(orgs))
	for i, org := range orgs {
		organizations[i] = org.Org
	}

	var flow *cache.FlowData

	if data.FlowID != nil {
		flow, err = cache.GetFlow(ctx, *data.FlowID)
		if err != nil {
			log.Info("Flow not found, creating a new one", zap.String("flowId", *data.FlowID), zap.Error(err))
			flow = nil
		}
	}

	if flow == nil {
		// If the flow is not found or is nil, create a new flow
		fId, err := uuid.NewV7()
		if err != nil {
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to generate flow ID!", http.StatusInternalServerError, err)
		}
		flow = &cache.FlowData{
			ID: fId.String(),
		}
	}

	flow.Type = cache.FlowTypeLogin
	flow.UserID = user.ID.String()
	flow.Email = user.Email
	flow.Orgs = organizations
	flow.MFARequired = false
	flow.MFAVerified = false
	flow.CreatedAt = time.Now()
	flow.ExpiresAt = time.Now().Add(10 * time.Minute)

	err = cache.StoreFlow(ctx, *flow)
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to store flow data!", http.StatusInternalServerError, err)
	}

	return &UserLoginResult{
		Message: "Multiple organizations found. Please select your organization to continue.",
		FlowID:  &flow.ID, // Return the flow ID for the user to continue the flow
	}, nil // Placeholder return
}
