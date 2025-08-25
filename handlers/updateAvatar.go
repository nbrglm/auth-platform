package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/nbrglm/auth-platform/config"
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

type UpdateAvatarHandler struct {
	AvatarCounterWEB *prometheus.CounterVec
	AvatarCounterAPI *prometheus.CounterVec
}

func NewUpdateAvatarHandlerWEB() *UpdateAvatarHandler {
	return &UpdateAvatarHandler{
		AvatarCounterWEB: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "update_avatar_web_requests",
				Help:      "Total number of update avatar web page requests",
			},
			[]string{"status"},
		),
		AvatarCounterAPI: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nbrglm_auth_platform",
				Subsystem: "auth",
				Name:      "update_avatar_api_requests",
				Help:      "Total number of update avatar API requests",
			},
			[]string{"status"},
		),
	}
}

func (h *UpdateAvatarHandler) Register(engine *gin.Engine) {
	metrics.Collectors = append(metrics.Collectors, h.AvatarCounterWEB, h.AvatarCounterAPI)

	// Register web routes
	engine.GET("/update-avatar", middlewares.RateLimitUIAuthenticatedMiddleware(), h.HandleUpdateAvatarGET)
	engine.POST("/update-avatar", middlewares.RateLimitUIAuthenticatedMiddleware(), h.HandleUpdateAvatarPOST)
	// engine.POST("/api/update-avatar", middlewares.RateLimitAPIMiddleware(), h.HandleUpdateAvatarAPI)
}

type UpdateAvatarPageParams struct {
	models.CommonPageParams
	CurrentAvatarURL string
}

func (h *UpdateAvatarHandler) HandleUpdateAvatarGET(c *gin.Context) {
	_, log, span := internal.WithContext(c.Request.Context(), "update_avatar_web_get")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	// Check if session exists
	if _, exists := c.Get(middlewares.CtxSessionExistsKey); !exists {
		c.Redirect(http.StatusFound, "/auth/s/refresh?returnTo=/update-avatar")
		return
	}

	claims, exist := c.Get(middlewares.CtxSessionTokenClaimsKey)
	if !exist {
		c.Redirect(http.StatusFound, "/auth/s/refresh?returnTo=/update-avatar")
		return
	}

	authClaims, ok := claims.(*tokens.AuthPlatformClaims)
	if !ok {
		log.Error("Invalid session claims type", zap.String("expected", "AuthPlatformClaims"), zap.Any("actual", claims), zap.String("operation", "update_avatar_web_get"))
		c.Redirect(http.StatusFound, "/auth/s/refresh?returnTo=/update-avatar")
		return
	}

	params := UpdateAvatarPageParams{
		CommonPageParams: models.NewCommonPageParams("Update Avatar", csrf.Token(c.Request)),
		CurrentAvatarURL: authClaims.UserAvatarURL,
	}

	c.HTML(http.StatusOK, "updateAvatar.html", params)
}

func (h *UpdateAvatarHandler) HandleUpdateAvatarPOST(c *gin.Context) {
	h.AvatarCounterWEB.WithLabelValues("received").Inc()

	ctx, log, span := internal.WithContext(c.Request.Context(), "update_avatar_web_post")
	defer span.End() // Ensure the span is ended to avoid memory leaks

	// Check if session exists
	if _, exists := c.Get(middlewares.CtxSessionExistsKey); !exists {
		c.Redirect(http.StatusFound, "/auth/s/refresh?returnTo=/update-avatar")
		return
	}

	claims, exist := c.Get(middlewares.CtxSessionTokenClaimsKey)
	if !exist {
		c.Redirect(http.StatusFound, "/auth/s/refresh?returnTo=/update-avatar")
		return
	}

	authClaims, ok := claims.(*tokens.AuthPlatformClaims)
	if !ok {
		log.Error("Invalid session claims type", zap.String("expected", "AuthPlatformClaims"), zap.Any("actual", claims), zap.String("operation", "update_avatar_web_post"))
		c.Redirect(http.StatusFound, "/auth/s/refresh?returnTo=/update-avatar")
		return
	}

	uid, err := uuid.Parse(authClaims.Subject)
	if err != nil {
		h.AvatarCounterWEB.WithLabelValues("error").Inc()
		log.Error("Invalid user ID in token claims", zap.Error(err), zap.String("operation", "update_avatar_web_post"))
		span.SetStatus(codes.Error, "Invalid user ID in token claims")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Please refresh the page and try again.", "invalid user id in token claims", http.StatusBadRequest, nil).Filter())
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		h.AvatarCounterWEB.WithLabelValues("error").Inc()
		log.Error("Failed to get avatar file from form", zap.Error(err), zap.String("operation", "update_avatar_web_post"))
		span.SetStatus(codes.Error, "Failed to get avatar file from form")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Please provide a valid avatar file.", "failed to get avatar file from form", http.StatusBadRequest, nil).Filter())
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/jpg" && contentType != "image/webp" {
		h.AvatarCounterWEB.WithLabelValues("invalid_request").Inc()
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Please upload a PNG, JPEG, or WEBP image.", "unsupported avatar file type", http.StatusBadRequest, nil).Filter())
		return
	}
	cacheControl := "public, max-age=" + strconv.Itoa(config.JWT.SessionTokenExpiration)

	key := "users/" + uid.String() + "/avatar"

	avatarKey, err := store.Objects.UploadPublicObject(ctx, key, file, contentType, cacheControl)
	if err != nil {
		h.AvatarCounterWEB.WithLabelValues("error").Inc()
		log.Error("Failed to upload avatar", zap.Error(err), zap.String("operation", "update_avatar_web_post"))
		span.SetStatus(codes.Error, "Failed to upload avatar")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("We encountered an error while updating your avatar. Please try again later.", "failed to upload avatar", http.StatusInternalServerError, err).Filter())
		return
	}

	url, err := store.Objects.GetObjectURL(ctx, avatarKey)
	if err != nil {
		h.AvatarCounterWEB.WithLabelValues("error").Inc()
		log.Error("Failed to get avatar URL", zap.Error(err), zap.String("operation", "update_avatar_web_post"))
		span.SetStatus(codes.Error, "Failed to get avatar URL")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("We encountered an error while updating your avatar. Please try again later.", "failed to get avatar url", http.StatusInternalServerError, nil).Filter())
		return
	}

	_, err = store.Querier.UpdateUser(ctx, db.UpdateUserParams{
		AvatarUrl: &url,
		Email:     authClaims.Email,
	})
	if err != nil {
		h.AvatarCounterWEB.WithLabelValues("error").Inc()
		log.Error("Failed to update user avatar URL in database", zap.Error(err), zap.String("operation", "update_avatar_web_post"))
		span.SetStatus(codes.Error, "Failed to update user avatar URL in database")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("We encountered an error while updating your avatar. Please try again later.", "failed to update user avatar url in database", http.StatusInternalServerError, nil).Filter())
		return
	}

	// Clear the session token so that new avatar URL is picked up
	tokens.RemoveTokens(c, true, false)

	h.AvatarCounterWEB.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Avatar updated successfully! It might take some time (%s minutes) to reflect everywhere.", strconv.Itoa(config.JWT.SessionTokenExpiration/60)),
	})
}
