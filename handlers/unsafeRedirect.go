package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/prometheus/client_golang/prometheus"
)

type UnsafeRedirectHandler struct {
	UnsafeRedirectCounter prometheus.Counter
}

func NewUnsafeRedirectHandler() *UnsafeRedirectHandler {
	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "nbrglm_auth_platform",
			Subsystem: "auth",
			Name:      "unsafe_redirect_requests",
			Help:      "Total number of unsafe redirect requests",
		},
	)
	return &UnsafeRedirectHandler{
		UnsafeRedirectCounter: counter,
	}
}

func (h *UnsafeRedirectHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.UnsafeRedirectCounter)

	engine.GET("/unsafe-redirect", middlewares.RateLimitUIOpenMiddleware(), h.HandleUnsafeRedirectGET)
}

type UnsafeRedirectPageParams struct {
	models.CommonPageParams
	RedirectURL *string
}

// HandleUnsafeRedirectGET handles the GET request for unsafe redirects.
func (h *UnsafeRedirectHandler) HandleUnsafeRedirectGET(c *gin.Context) {
	h.UnsafeRedirectCounter.Inc()

	var redirectURL *string
	if v := c.Query("original"); v != "" {
		redirectURL = &v
	}

	c.HTML(http.StatusOK, "unsafeRedirect.html", UnsafeRedirectPageParams{
		CommonPageParams: models.NewCommonPageParams("Unsafe Redirection", csrf.Token(c.Request)),
		RedirectURL:      redirectURL,
	})
}
