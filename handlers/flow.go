package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/nexeres/internal/cache"
	"github.com/nbrglm/nexeres/internal/metrics"
	"github.com/nbrglm/nexeres/internal/models"
	"github.com/prometheus/client_golang/prometheus"
)

type FlowHandler struct {
	FlowGetCounter *prometheus.CounterVec
}

func NewFlowHandler() *FlowHandler {
	return &FlowHandler{
		FlowGetCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexeres",
				Subsystem: "auth",
				Name:      "user_flow_get_requests",
				Help:      "Total number of get requests for user flow data",
			},
			[]string{"status"},
		),
	}
}

func (h *FlowHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.FlowGetCounter)
	engine.GET("/api/auth/flow/:flowId", h.HandleGetFlow)
}

type UserFlowData cache.FlowData

// HandleGetFlow godoc
// @Summary Get User Flow Data
// @Description Retrieves user flow data based on the provided flow ID.
// @Tags Auth
// @Accept json
// @Produce json
// @Param flowId path string true "Flow ID"
// @Success 200 {object} UserFlowData "User Flow Data"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/auth/flow/{flowId} [get]
func (h *FlowHandler) HandleGetFlow(c *gin.Context) {
	h.FlowGetCounter.WithLabelValues("received").Inc()

	flowId := strings.TrimSuffix(strings.TrimSpace(c.Param("flowId")), "/")

	flow, err := cache.GetFlow(c.Request.Context(), flowId)

	if err != nil {
		if err == cache.ErrKeyNotFound {
			h.FlowGetCounter.WithLabelValues("not_found").Inc()
			c.JSON(http.StatusBadRequest, models.NewErrorResponse("Flow not found", "Not Found", http.StatusBadRequest, nil))
			return
		}
		h.FlowGetCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.GenericErrorMessage, "Error fetching flow data", http.StatusBadRequest, nil))
		return
	}

	h.FlowGetCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, flow)
}
