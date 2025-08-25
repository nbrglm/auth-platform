package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/nbrglm/auth-platform/db"
	"github.com/nbrglm/auth-platform/internal"
	"github.com/nbrglm/auth-platform/internal/metrics"
	"github.com/nbrglm/auth-platform/internal/middlewares"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/store"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type RootHandlerWEB struct {
	ProfileCounter *prometheus.CounterVec
}

func NewRootHandlerWEB() *RootHandlerWEB {
	return &RootHandlerWEB{
		ProfileCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "root_web_requests",
				Help:      "Total number of root web page requests",
			},
			[]string{"status"},
		),
	}
}

func (h *RootHandlerWEB) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.ProfileCounter)
	engine.GET("/", middlewares.RateLimitUIAuthenticatedMiddleware(), h.HandleRootGET)
	// engine.POST("/profile/update", h.HandleProfileUpdatePOST) // Uncomment when implemented
}

type RootProfilePageParams struct {
	models.CommonPageParams
	User *db.GetUserByIDRow
}

func (h *RootHandlerWEB) HandleRootGET(c *gin.Context) {
	// Increment the counter for root page GET requests
	h.ProfileCounter.WithLabelValues("success").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "root_profile_web_get")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	// Check if session exists
	if _, exists := c.Get(middlewares.CtxSessionExistsKey); !exists {
		// If no session, redirect to session refresh
		c.Redirect(http.StatusSeeOther, "/auth/s/refresh?returnTo=/")
		return
	}

	claims, exist := c.Get(middlewares.CtxSessionTokenClaimsKey)
	if !exist {
		// If no claims, redirect to session refresh
		c.Redirect(http.StatusSeeOther, "/auth/s/refresh?returnTo=/")
		return
	}

	var authClaims *tokens.AuthPlatformClaims
	var ok bool

	if authClaims, ok = claims.(*tokens.AuthPlatformClaims); !ok {
		// If claims are not of expected type, redirect to session refresh
		c.Redirect(http.StatusSeeOther, "/auth/s/refresh?returnTo=/")
		return
	}

	uid, err := uuid.Parse(authClaims.Subject)
	if err != nil {
		h.ProfileCounter.WithLabelValues("error").Inc()
		log.Error("invalid user ID in token claims", zap.Error(err))
		span.SetStatus(codes.Error, "invalid user ID in token claims")
		// If user ID is invalid, redirect to session refresh
		c.Redirect(http.StatusSeeOther, "/auth/s/refresh?returnTo=/")
		return
	}

	user, err := store.Querier.GetUserByID(ctx, uid)
	if err != nil {
		h.ProfileCounter.WithLabelValues("error").Inc()
		var params models.CommonPageParams
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("no user found for id obtained from auth claims", zap.Error(err), zap.String("operation", "root_profile_web_get"))
			span.SetStatus(codes.Error, "no user found for id obtained from auth claims")
			params = models.NewPageError("No such user!", "It seems no user exists. This might be due to your account being deleted OR you being banned from the platform. It is recommended you logout, and try logging in again.", "/auth/s/logout", "Logout")
		} else {
			log.Error("failed to get user by ID", zap.Error(err), zap.String("operation", "root_profile_web_get"))
			span.SetStatus(codes.Error, "failed to get user by ID")
			params = models.NewPageError("An Error Occurred!", "We encountered an error while fetching your profile. Please try again later.", "/auth/s/refresh?returnTo=/", "Retry")
		}
		c.HTML(http.StatusOK, "root.html", params)
		return
	}

	// Render the root page
	c.HTML(http.StatusOK, "root.html", RootProfilePageParams{
		CommonPageParams: models.NewCommonPageParams("Profile", csrf.Token(c.Request)),
		User:             &user,
	})
}
