package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/nbrglm/auth-platform/internal"
	"github.com/nbrglm/auth-platform/internal/auth"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/prometheus/client_golang/prometheus"
)

type VerifyEmailWEBHandler struct {
	VerifyEmailGetCounter  *prometheus.CounterVec
	VerifyEmailPostCounter *prometheus.CounterVec
}

// NewVerifyEmailWEBHandler creates a new VerifyEmailWEBHandler instance.
func NewVerifyEmailWEBHandler() *VerifyEmailWEBHandler {
	return &VerifyEmailWEBHandler{
		VerifyEmailGetCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_verify_email_web_get_requests",
				Help:      "Total number of user email verification GET requests",
			},
			[]string{"status"},
		),
		VerifyEmailPostCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_verify_email_web_post_requests",
				Help:      "Total number of user email verification POST requests",
			},
			[]string{"status"},
		),
	}
}

func (h *VerifyEmailWEBHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.VerifyEmailGetCounter)
	metrics.Collectors = append(metrics.Collectors, h.VerifyEmailPostCounter)

	engine.GET("/auth/verify-email", h.HandleVerifyEmailGET)
	engine.POST("/auth/verify-email/submit", h.HandleVerifyEmailPOST)
	engine.GET("/auth/verify-email/error", h.HandleVerifyEmailErrorGET)
	engine.GET("/auth/verify-email/success", h.HandleVerifyEmailSuccessGET)
}

// VerifyEmailPageParams defines the parameters for the verify email page.
//
// For displaying an error, pass a models.PageError to c.HTML.
type VerifyEmailPageParams struct {
	models.CommonPageParams

	// If a verification token is provided, a post form will be rendered.
	// Otherwise, a normal 'email is sent' page will be rendered.
	VerificationToken *string

	// Success denotes whether the email verification was successful.
	// If so, it renderes the success page.
	Success bool

	// If a non-default success link is provided, it will be used to redirect the user after successful verification.
	SuccessLink *string
}

func (h *VerifyEmailWEBHandler) HandleVerifyEmailGET(c *gin.Context) {
	h.VerifyEmailGetCounter.WithLabelValues("received").Inc()

	var token *string
	t := strings.TrimSpace(c.Query("token"))
	if t != "" {
		token = &t
	}

	h.VerifyEmailGetCounter.WithLabelValues("success").Inc()

	// Render the verification email page
	c.HTML(http.StatusOK, "auth/verify_email.html", VerifyEmailPageParams{
		CommonPageParams:  models.NewCommonPageParams("Verify Email", csrf.Token(c.Request)),
		VerificationToken: token,
	})
}

func (h *VerifyEmailWEBHandler) HandleVerifyEmailPOST(c *gin.Context) {
	h.VerifyEmailPostCounter.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "verify_email_token_web")
	defer span.End()

	var data auth.VerifyEmailTokenData
	if err := c.ShouldBind(&data); err != nil {
		h.VerifyEmailPostCounter.WithLabelValues("invalid_input").Inc()
		setErrorAndRedirect(c, models.NewErrorResponse("Invalid request! You may have followed an invalid link. Please try again, by clicking on the retry verification button. You will be redirected to the login page, once you fill in your details, your new verification email will be sent.", "Failed to bind login data", http.StatusBadRequest, nil).WithUI("/auth/login", "/auth/verify-email/error", "Retry Verification"))
		return
	}

	data.Token = strings.TrimSpace(data.Token)
	if data.Token == "" {
		h.VerifyEmailPostCounter.WithLabelValues("missing_token").Inc()
		setErrorAndRedirect(c, models.NewErrorResponse("Missing verification token! Please try again. You will be redirected to the login page, once you fill in your details, your new verification email will be sent.", "Failed to bind login data", http.StatusBadRequest, nil).WithUI("/auth/login", "/auth/verify-email/error", "Retry Verification"))
		return
	}

	response, err := auth.VerifyEmailToken(ctx, log, data)
	response = ProcessUiResult(c, response, err.WithUI("/auth/login", "/auth/verify-email/error", "Retry Verification"), span, log, h.VerifyEmailPostCounter, "verify_email_token_web")

	if response == nil {
		return
	}

	h.VerifyEmailPostCounter.WithLabelValues("success").Inc()
	// Redirect to the success page
	c.Redirect(http.StatusSeeOther, "/auth/verify-email/success")
}

func (h *VerifyEmailWEBHandler) HandleVerifyEmailErrorGET(c *gin.Context) {
	// No need to increment a counter here, as this is an error page.

	params := getPageError(c, "Email Verification Error")

	// Render the error page
	c.HTML(http.StatusOK, "auth/verify_email.html", params)
}

func (h *VerifyEmailWEBHandler) HandleVerifyEmailSuccessGET(c *gin.Context) {
	// No need to increment a counter here, as this is a success page.

	successLink := "/auth/login"
	// Render the success page
	c.HTML(http.StatusOK, "auth/verify_email.html", VerifyEmailPageParams{
		CommonPageParams: models.NewCommonPageParams("Email Verified Successfully", csrf.Token(c.Request)),
		Success:          true,
		SuccessLink:      &successLink,
	})
}
