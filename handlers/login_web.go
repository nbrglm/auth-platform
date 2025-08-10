package handlers

import (
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/nbrglm/auth-platform/internal"
	"github.com/nbrglm/auth-platform/internal/auth"
	"github.com/nbrglm/auth-platform/internal/cache"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"github.com/prometheus/client_golang/prometheus"
)

type LoginWEBHandler struct {
	LoginGetCounter  *prometheus.CounterVec
	LoginPostCounter *prometheus.CounterVec
}

func NewLoginWEBHandler() *LoginWEBHandler {
	return &LoginWEBHandler{
		LoginGetCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_login_web_get_requests",
				Help:      "Total number of user login GET requests",
			},
			[]string{"status"},
		),
		LoginPostCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_login_web_post_requests",
				Help:      "Total number of user login POST requests",
			},
			[]string{"status"},
		),
	}
}

func (h *LoginWEBHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.LoginGetCounter)
	metrics.Collectors = append(metrics.Collectors, h.LoginPostCounter)

	engine.GET("/auth/login", middlewares.RateLimitUIOpenMiddleware(), middlewares.NewRedirectIfAuthenticatedMiddleware(), h.HandleLoginGET)
	engine.POST("/auth/login/submit", middlewares.RateLimitUIOpenMiddleware(), h.HandleLoginPOST)
	engine.GET("/auth/login/error", middlewares.RateLimitUIOpenMiddleware(), h.HandleLoginErrorGET)
}

type LoginPageParams struct {
	models.CommonPageParams
	FlowID string
}

// HandleLoginGET handles the GET request for the login page.
func (h *LoginWEBHandler) HandleLoginGET(c *gin.Context) {
	h.LoginGetCounter.WithLabelValues("received").Inc()

	// Get the return URL from the context, if it exists
	returnTo := c.GetString(middlewares.ReturnToURLKey)
	unsafe := c.GetBool(middlewares.ReturnToURLUnsafeKey)
	retryUrl := "/auth/login?returnTo="
	if returnTo == "" {
		// Default to a safe redirect URL if none is provided
		// and show the page with go to home button
		returnTo = "/"

		// No returnTo URL, so we set a default retry URL
		retryUrl = "/auth/login"
	} else if unsafe {
		// Add the original returnTo URL as a query parameter for retryUrl
		retryUrl += url.QueryEscape(returnTo)
		returnTo = "/unsafe-redirect?original=" + url.QueryEscape(returnTo)
	} else {
		// If returnTo is provided and safe, we use it directly
		retryUrl += url.QueryEscape(returnTo)
	}

	flowId, err := uuid.NewV7()
	if err != nil {
		h.LoginGetCounter.WithLabelValues("error").Inc()
		setErrorAndRedirect(c, models.NewErrorResponse(models.GenericErrorMessage, "Could not generate flow id!", http.StatusInternalServerError, nil).WithUI(retryUrl, "/auth/login/error", "Retry Login"))
		return
	}

	err = cache.StoreFlow(c.Request.Context(), cache.FlowData{
		ID:       flowId.String(),
		Type:     cache.FlowTypeLogin,
		ReturnTo: returnTo,
	})

	if err != nil {
		h.LoginGetCounter.WithLabelValues("error").Inc()
		setErrorAndRedirect(c, models.NewErrorResponse(models.GenericErrorMessage, "Could not store flow data!", http.StatusInternalServerError, err).WithUI(retryUrl, "/auth/login/error", "Retry Login"))
		return
	}

	h.LoginGetCounter.WithLabelValues("success").Inc()
	c.HTML(http.StatusOK, "auth/login.html", LoginPageParams{
		CommonPageParams: models.NewCommonPageParams("Login", csrf.Token(c.Request)),
		FlowID:           flowId.String(),
	})
}

// HandleLoginPOST handles the POST request for user login.
func (h *LoginWEBHandler) HandleLoginPOST(c *gin.Context) {
	h.LoginPostCounter.WithLabelValues("received").Inc()

	var data auth.UserLoginData
	if err := c.ShouldBind(&data); err != nil {
		h.LoginPostCounter.WithLabelValues("invalid_input").Inc()
		setErrorAndRedirect(c, models.NewErrorResponse("Invalid input data", "Failed to bind login data", http.StatusBadRequest, nil).WithUI("/auth/login", "/auth/login/error", "Retry Login"))
		return
	}

	ctx, log, span := internal.WithContext(c.Request.Context(), "login_web")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	returnTo := "/unsafe-redirect"

	// Get the return URL from the flowId
	if data.FlowID != nil {
		flow, err := cache.GetFlow(ctx, *data.FlowID)
		if err == nil && flow != nil {
			if strings.TrimSpace(flow.ReturnTo) != "" {
				returnTo = flow.ReturnTo
			}
		}
	}

	data.IPAddress = netip.MustParseAddr(c.ClientIP())
	data.UserAgent = c.Request.UserAgent()

	response, err := auth.HandleLogin(ctx, log, data)
	response = ProcessUiResult(c, response, err.WithUI("/auth/login?returnTo="+url.QueryEscape(returnTo), "/auth/login/error", "Retry Login"), span, log, h.LoginPostCounter, "login_web")
	if response == nil {
		return
	}

	if response.SentVerificationEmail {
		// Redirect to the verification page if the email is not verified
		h.LoginPostCounter.WithLabelValues("email_not_verified").Inc()
		c.Redirect(http.StatusSeeOther, "/auth/verify-email")
		return
	}

	// Set the tokens in cookies
	tks := response.Tokens
	tokens.SetTokens(c, &tks)

	h.LoginPostCounter.WithLabelValues("success").Inc()
	c.Redirect(http.StatusSeeOther, returnTo)
}

func (h *LoginWEBHandler) HandleLoginErrorGET(c *gin.Context) {
	params := getPageError(c, "Login Error")

	c.HTML(http.StatusOK, "auth/login.html", params)
}
