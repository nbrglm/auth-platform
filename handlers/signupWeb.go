package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/nbrglm/auth-platform/internal"
	"github.com/nbrglm/auth-platform/internal/auth"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/prometheus/client_golang/prometheus"
)

type SignupWEBHandler struct {
	SignupGetCounter  *prometheus.CounterVec
	SignupPostCounter *prometheus.CounterVec
}

// NewSignupWEBHandler creates a new SignupWEBHandler instance.
func NewSignupWEBHandler() *SignupWEBHandler {
	return &SignupWEBHandler{
		SignupGetCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_signup_web_get_requests",
				Help:      "Total number of user signup GET requests",
			},
			[]string{"status"},
		),
		SignupPostCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_signup_web_post_requests",
				Help:      "Total number of user signup POST requests",
			},
			[]string{"status"},
		),
	}
}

func (h *SignupWEBHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.SignupGetCounter)
	metrics.Collectors = append(metrics.Collectors, h.SignupPostCounter)

	engine.GET("/auth/signup", middlewares.RateLimitUIOpenMiddleware(), h.HandleSignupGET)
	engine.POST("/auth/signup/submit", middlewares.RateLimitUIOpenMiddleware(), h.HandleSignupPOST)
	engine.GET("/auth/signup/error", middlewares.RateLimitUIOpenMiddleware(), h.HandleSignupErrorGET)
}

// HandleSignupGET handles the GET request for the signup page.
func (h *SignupWEBHandler) HandleSignupGET(c *gin.Context) {
	h.SignupGetCounter.WithLabelValues("received").Inc()

	params := models.NewCommonPageParams("Sign Up", csrf.Token(c.Request))

	c.HTML(http.StatusOK, "auth/signup.html", params)
}

func (h *SignupWEBHandler) HandleSignupPOST(c *gin.Context) {
	h.SignupPostCounter.WithLabelValues("received").Inc()

	var data auth.UserSignupData
	if err := c.ShouldBind(&data); err != nil {
		h.SignupPostCounter.WithLabelValues("invalid_input").Inc()
		var errResponse *models.ErrorResponse
		if isUnmatchingPasswordsError(err) {
			errResponse = models.NewErrorResponse("Passwords do not match! Please try again.", "", http.StatusSeeOther, nil)
		} else {
			errResponse = models.NewErrorResponse("Invalid input data. Please check your input and try again.", "", http.StatusSeeOther, nil)
		}
		setErrorAndRedirect(c, errResponse.WithUI("/auth/signup", "/auth/signup/error", "Retry Signup"))
		return
	}

	ctx, log, span := internal.WithContext(c.Request.Context(), "signup_web")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	// Call the internal handler to process the signup
	response, err := auth.HandleSignup(ctx, log, data)

	// Process Errors and Edge Cases
	response = ProcessUiResult(c, response, err.WithUI("/auth/signup", "/auth/signup/error", "Retry Signup"), span, log, h.SignupPostCounter, "Sign Up")

	if response == nil {
		// A nil response indicates an error occurred, which has already been handled in ProcessUiResult.
		// We can return early here as the error handling logic has already set the appropriate session data and redirected the user.
		return
	}

	// Send the verification email since the signup was successful.
	emailResponse, err := auth.SendVerificationEmail(ctx, log, auth.SendVerificationEmailData{
		Email: data.Email,
	})

	if err != nil || emailResponse == nil || !emailResponse.Success {
		if err != nil {
			// Change the error message to indicate that the user has been signed up successfully
			err.Message = "You have been signed up successfully, but an error occurred while sending the verification email. Please continue to login, which will send you a verification email."
			ProcessUiResult(c, emailResponse, err.WithUI("/auth/login", "/auth/signup/error", "Login"), span, log, h.SignupPostCounter, "Sign Up")
		} else {
			// If it's not an ErrorResponse, we create a new one with a generic message
			errResponse := models.NewErrorResponse("You have been signed up successfully, but an error occurred while sending the verification email. Please continue to login, which will send you a verification email.", "Failed to send verification email!", http.StatusInternalServerError, err)
			ProcessUiResult(c, emailResponse, errResponse.WithUI("/auth/login", "/auth/signup/error", "Login"), span, log, h.SignupPostCounter, "Sign Up")
		}
		return
	}

	// If the signup was successful, we redirect to the home page
	h.SignupPostCounter.WithLabelValues("success").Inc()
	c.Redirect(http.StatusSeeOther, "/auth/verify-email")
}

func isUnmatchingPasswordsError(err error) bool {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		for _, fieldErr := range validationErr {
			if fieldErr.Field() == "ConfirmPassword" && fieldErr.Tag() == "eqfield" {
				return true
			}
		}
	}
	return false
}

func (h *SignupWEBHandler) HandleSignupErrorGET(c *gin.Context) {
	params := getPageError(c, "Signup Error")

	c.HTML(http.StatusOK, "auth/signup.html", params)
}
