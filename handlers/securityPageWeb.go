package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/prometheus/client_golang/prometheus"
)

type SecurityPageHandler struct {
	SecurityPageCounter      *prometheus.CounterVec
	UpdatePasswordCounterWEB *prometheus.CounterVec
	UpdatePasswordCounterAPI *prometheus.CounterVec
}

func NewSecurityPageHandler() *SecurityPageHandler {
	return &SecurityPageHandler{
		SecurityPageCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "security_page_web_requests",
				Help:      "Total number of security page web requests",
			},
			[]string{"status"},
		),
		UpdatePasswordCounterWEB: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "update_password_web_requests",
				Help:      "Total number of update password web requests",
			},
			[]string{"status"},
		),
		UpdatePasswordCounterAPI: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "update_password_api_requests",
				Help:      "Total number of update password API requests",
			},
			[]string{"status"},
		),
	}
}

func (h *SecurityPageHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.SecurityPageCounter, h.UpdatePasswordCounterWEB, h.UpdatePasswordCounterAPI)
	engine.GET("/security", middlewares.RateLimitUIAuthenticatedMiddleware(), h.HandleSecurityGET)
	// engine.POST("/security/update-password", middlewares.RateLimitUIAuthenticatedMiddleware(), h.HandleUpdatePasswordPOST)
}

func (h *SecurityPageHandler) HandleSecurityGET(c *gin.Context) {}
