package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type LogoutHandler struct {
	LogoutWEBCounter *prometheus.CounterVec
	LogoutAPICounter *prometheus.CounterVec
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{
		LogoutWEBCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_logout_web_requests",
				Help:      "Total number of user WEB logout requests",
			},
			[]string{"status"},
		),
		LogoutAPICounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "user_logout_api_requests",
				Help:      "Total number of user API logout requests",
			},
			[]string{"status"},
		),
	}
}

func (h *LogoutHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.LogoutWEBCounter)
	metrics.Collectors = append(metrics.Collectors, h.LogoutAPICounter)

	// engine.GET("/auth/s/logout", h.HandleLogoutWEB)
	// engine.POST("/api/auth/s/logout", h.HandleLogoutAPI)
}
