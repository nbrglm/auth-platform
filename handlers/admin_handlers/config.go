package admin_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/nexeres/config"
	"github.com/nbrglm/nexeres/internal/metrics"
	"github.com/nbrglm/nexeres/internal/middlewares"
	"github.com/prometheus/client_golang/prometheus"
)

type ConfigHandler struct {
	ConfigGETCounter *prometheus.CounterVec
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		ConfigGETCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "admin",
				Name:      "config_get_requests_total",
				Help:      "Total number of admin configuration read requests",
			},
			[]string{"status"},
		),
	}
}

func (h *ConfigHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.ConfigGETCounter)
	engine.GET("/api/admin/config", middlewares.RequireAuth(middlewares.AuthModeAdmin), h.GetConfig)
}

// GetConfig godoc
// @Summary      Get current configuration
// @Description  Retrieves the current configuration of the application.
// @Tags         Admin
// @Produce      json
// @Success 200 {object}  config.CompleteConfig "Current configuration"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/admin/config [get]
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	h.ConfigGETCounter.WithLabelValues("received").Inc()

	h.ConfigGETCounter.WithLabelValues("success").Inc()
	middlewares.AdminInactivityReset(c) // Reset inactivity timer
	c.JSON(http.StatusOK, config.Config)
}
