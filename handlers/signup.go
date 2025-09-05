package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nbrglm/nexeres/config"
	"github.com/nbrglm/nexeres/db"
	"github.com/nbrglm/nexeres/internal"
	"github.com/nbrglm/nexeres/internal/metrics"
	"github.com/nbrglm/nexeres/internal/models"
	"github.com/nbrglm/nexeres/internal/password"
	"github.com/nbrglm/nexeres/internal/store"
	"github.com/nbrglm/nexeres/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type SignupHandler struct {
	SignupCounter *prometheus.CounterVec
}

func NewSignupHandler() *SignupHandler {
	return &SignupHandler{
		SignupCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "auth",
				Name:      "user_signup_requests",
				Help:      "Total number of user signup requests",
			},
			[]string{"status"},
		),
	}
}

func (h *SignupHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.SignupCounter)
	engine.POST("/api/auth/signup", h.HandleSignup)
}

type UserSignupData struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	FirstName       string `json:"firstName" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	InviteToken     string `json:"inviteToken,omitempty"` // Optional invite token for signup
}

type UserSignupResult struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

// HandleSignup godoc
// @Summary User Signup
// @Description Handles user registration requests.
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body UserSignupData true "User Signup Data"
// @Success 200 {object} UserSignupResult "User Signup Result"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Invalid Invite Token or Missing Invite Token or Domain Not Allowed"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/signup [post]
func (h *SignupHandler) HandleSignup(c *gin.Context) {
	h.SignupCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "signup")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	var signupData UserSignupData
	if err := c.ShouldBindJSON(&signupData); err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Invalid request data", "Please check your input and try again.", http.StatusBadRequest, nil), span, log, h.SignupCounter, "signup")
		return
	}

	domain, err := utils.GetDomainFromEmail(signupData.Email)
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse("Invalid request! Please input a valid email and try again.", "Invalid email domain!", http.StatusBadRequest, nil), span, log, h.SignupCounter, "signup")
		return
	}

	// If multitenancy is enabled, we will use the invite token to determine the organization,
	// if not present, then we will check the domain against the existing organizations,
	// if the domain matches a valid, auto-join enabled, verified existing domain, we will use that organization.
	// Otherwise if multitenancy is not enabled, we will fetch the default organization.

	var org *db.Org
	role := "member" // Default role for new users

	if config.Multitenancy {
		if strings.TrimSpace(signupData.InviteToken) == "" {
			// Check if the domain matches any verified, auto-join enabled domains for any organization
			organization, err := store.Querier.GetOrgForDomainIfAutoJoin(ctx, domain)
			if errors.Is(err, pgx.ErrNoRows) {
				utils.ProcessError(c, models.NewErrorResponse("Email is not associated to any Organizations! Please contact your administrator.", "No organization found for domain that has auto-join enabled! If you have not verified the domain yet, please do so.", http.StatusUnauthorized, nil), span, log, h.SignupCounter, "signup")
				return
			}

			if err != nil {
				utils.ProcessError(c, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to get organization for domain!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
				return
			}
			// we do not change role here, as we are not using an invite token,
			// so the user will be a member by default.
			org = &organization
		} else {
			invitation, err := store.Querier.GetInvitationByToken(ctx, signupData.InviteToken)
			if errors.Is(err, pgx.ErrNoRows) {
				utils.ProcessError(c, models.NewErrorResponse("Invalid invite token! Please check your token and try again.", "No invitation found for the provided token!", http.StatusUnauthorized, nil), span, log, h.SignupCounter, "signup")
				return
			}
			if err != nil {
				utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to get invitation by token!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
				return
			}
			if invitation.Email != signupData.Email {
				utils.ProcessError(c, models.NewErrorResponse("Invalid invite token! Please check your token and try again.", "The invite token does not match the provided email!", http.StatusUnauthorized, nil), span, log, h.SignupCounter, "signup")
				return
			}

			// The invite is valid, let's get the organization and role from the invitation
			organization, err := store.Querier.GetOrgByID(ctx, invitation.OrgID)
			if errors.Is(err, pgx.ErrNoRows) {
				utils.ProcessError(c, models.NewErrorResponse("Invalid invite token! Please check your token and try again.", "No organization found for the provided invitation!", http.StatusUnauthorized, nil), span, log, h.SignupCounter, "signup")
				return
			}
			if err != nil {
				utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to get organization by ID for given invite token!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
				return
			}
			role = invitation.Role // Use the role from the invitation
			org = &organization
		}
	} else {
		organization, err := store.Querier.GetOrgBySlug(ctx, "default")
		if err != nil {
			// we return the underlying error here because default org is ALWAYS supposed to be found
			utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to get default organization!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
			return
		}
		org = &organization
	}

	id, err := uuid.NewV7()
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to generate user ID!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
		return
	}

	passwordHash, err := password.HashPassword(signupData.Password)
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to hash password!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
		return
	}

	tx, err := store.PgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to begin transaction!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
		return
	}
	defer tx.Rollback(ctx) // Ensure the transaction is rolled back if not committed

	q := store.Querier.WithTx(tx)

	// Create the user
	user, err := q.CreateUser(ctx, db.CreateUserParams{
		ID:           id,
		Email:        signupData.Email,
		PasswordHash: &passwordHash,
		FirstName:    &signupData.FirstName,
		LastName:     &signupData.LastName,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			utils.ProcessError(c, models.NewErrorResponse("Email is already registered! Please login or use a different email.", "User with this email already exists!", http.StatusBadRequest, nil), span, log, h.SignupCounter, "signup")
			return
		}
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to create user!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
		return
	}

	// Create the user organization membership
	if err := q.LinkUserToOrg(ctx, db.LinkUserToOrgParams{
		UserID: user.ID,
		OrgID:  org.ID,
		Role:   role,
	}); err != nil {
		utils.ProcessError(c, models.NewErrorResponse(models.GenericErrorMessage, "Failed to link user to organization!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
		return
	}

	// Commit the transaction, user creating is successful!
	if err := tx.Commit(ctx); err != nil {
		utils.ProcessError(c, models.NewErrorResponse("An error occurred while processing your request. Please try again later.", "Failed to commit transaction!", http.StatusInternalServerError, err), span, log, h.SignupCounter, "signup")
		return
	}

	// Implementation of the signup logic goes here.
	// This is a placeholder to illustrate where the actual signup handling code would be placed.
	c.JSON(http.StatusOK, &UserSignupResult{
		UserID:  user.ID.String(),
		Message: "Signup successful!",
	})
}
