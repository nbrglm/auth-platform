package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/db"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/password"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/utils"
	"go.uber.org/zap"
)

type UserSignupData struct {
	Email           string `json:"email" form:"email" binding:"required,email"`
	Password        string `json:"password" form:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" binding:"required,eqfield=Password"`
	FirstName       string `json:"firstName" form:"firstName" binding:"required"`
	LastName        string `json:"lastName" form:"lastName" binding:"required"`
	InviteToken     string `json:"inviteToken,omitempty" form:"inviteToken,omitempty"` // Optional invite token for signup
}

type UserSignupResult struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

func HandleSignup(ctx context.Context, log *zap.Logger, data UserSignupData) (*UserSignupResult, *models.ErrorResponse) {
	domain, err := utils.GetDomainFromEmail(data.Email)
	if err != nil {
		return nil, models.NewErrorResponse("Invalid request! Please input a valid email and try again.", "Invalid email domain!", http.StatusBadRequest, nil)
	}

	// If multitenancy is enabled, we will use the invite token to determine the organization,
	// if not present, then we will check the domain against the existing organizations,
	// if the domain matches a valid, auto-join enabled, verified existing domain, we will use that organization.
	// Otherwise if multitenancy is not enabled, we will fetch the default organization.

	var org *db.Org
	role := "member" // Default role for new users

	if config.Multitenancy.Enable {
		if strings.TrimSpace(data.InviteToken) == "" {
			// Check if the domain matches any verified, auto-join enabled domains for any organization
			organization, err := store.Querier.GetOrgForDomainIfAutoJoin(ctx, domain)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, models.NewErrorResponse("Email is not associated to any Organizations! Please contact your administrator.", "No organization found for domain that has auto-join enabled! If you have not verified the domain yet, please do so.", http.StatusBadRequest, nil)
			}

			if err != nil {
				return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to get organization for domain!", http.StatusInternalServerError, err)
			}
			// we do not change role here, as we are not using an invite token,
			// so the user will be a member by default.
			org = &organization
		} else {
			invitation, err := store.Querier.GetInvitationByToken(ctx, data.InviteToken)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, models.NewErrorResponse("Invalid invite token! Please check your token and try again.", "No invitation found for the provided token!", http.StatusBadRequest, nil)
			}
			if err != nil {
				return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to get invitation by token!", http.StatusInternalServerError, err)
			}
			if invitation.Email != data.Email {
				return nil, models.NewErrorResponse("Invalid invite token! Please check your token and try again.", "The invite token does not match the provided email!", http.StatusBadRequest, nil)
			}

			// The invite is valid, let's get the organization and role from the invitation
			organization, err := store.Querier.GetOrgByID(ctx, invitation.OrgID)
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, models.NewErrorResponse("Invalid invite token! Please check your token and try again.", "No organization found for the provided invitation!", http.StatusBadRequest, nil)
			}
			if err != nil {
				return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to get organization by ID for given invite token!", http.StatusInternalServerError, err)
			}
			role = invitation.Role // Use the role from the invitation
			org = &organization
		}
	} else {
		organization, err := store.Querier.GetOrgBySlug(ctx, "default")
		if err != nil {
			// we return the underlying error here because default org is ALWAYS supposed to be found
			return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to get default organization!", http.StatusInternalServerError, err)
		}
		org = &organization
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to generate user ID!", http.StatusInternalServerError, err)
	}

	passwordHash, err := password.HashPassword(data.Password)
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to hash password!", http.StatusInternalServerError, err)
	}

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to begin transaction!", http.StatusInternalServerError, err)
	}
	defer tx.Rollback(ctx) // Ensure the transaction is rolled back if not committed

	q := store.Querier.WithTx(tx)

	// Create the user
	user, err := q.CreateUser(ctx, db.CreateUserParams{
		ID:           id,
		Email:        data.Email,
		PasswordHash: &passwordHash,
		FirstName:    &data.FirstName,
		LastName:     &data.LastName,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, models.NewErrorResponse("Email already exists! Please use a different email.", "User with the provided email already exists!", http.StatusConflict, nil)
		}
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to create user!", http.StatusInternalServerError, err)
	}

	// Create the user organization membership
	if err := q.LinkUserToOrg(ctx, db.LinkUserToOrgParams{
		UserID: user.ID,
		OrgID:  org.ID,
		Role:   role,
	}); err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to link user to organization!", http.StatusInternalServerError, err)
	}

	// Commit the transaction, user creating is successful!
	if err := tx.Commit(ctx); err != nil {
		return nil, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to commit transaction!", http.StatusInternalServerError, err)
	}

	return &UserSignupResult{
		UserID:  user.ID.String(),
		Message: "User created successfully! Please verify your email to continue.",
	}, nil
}
