package cache

import (
	"context"
	"fmt"
	"time"

	gocache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	cache_metrics "github.com/eko/gocache/lib/v4/metrics"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/opts"
	"github.com/redis/go-redis/v9"
)

var cached *marshaler.Marshaler

func InitCache() error {
	redisPassword := ""
	if config.Stores.Redis.Password != nil {
		redisPassword = *config.Stores.Redis.Password
	}
	redisOpts := redis.Options{
		Addr:     config.Stores.Redis.Address,
		Password: redisPassword,
		DB:       config.Stores.Redis.DB,
	}
	redisClient := redis.NewClient(&redisOpts)

	// Initialize the Redis cache
	redisStore := redis_store.NewRedis(redisClient)

	metrics := cache_metrics.NewPrometheus(fmt.Sprintf("%s_cache", opts.Name))

	cacheManager := gocache.NewMetric(metrics, redisStore)
	cached = marshaler.New(cacheManager)
	return nil
}

type FlowType string

var (
	FlowTypeLogin         FlowType = "login"          // For Login Flow
	FlowTypeInvite        FlowType = "invite"         // For Invite Flow
	FlowTypePasswordReset FlowType = "password_reset" // For password reset flow
	FlowTypeSSO           FlowType = "sso"            // TODO: Implement SSO flow
)

type FlowData struct {
	ID          string             `json:"id"`
	Type        FlowType           `json:"type"`
	UserID      string             `json:"userId"`
	Email       string             `json:"email"`
	Orgs        []models.OrgCompat `json:"orgs,omitempty"`
	MFARequired bool               `json:"mfaRequired"`
	MFAVerified bool               `json:"mfaVerified"`
	InviteToken string             `json:"inviteToken,omitempty"` // For Invite Flow
	SSOProvider string             `json:"ssoProvider,omitempty"` // For SSO Flow, e.g., "google", "github", etc.
	SSOUserID   string             `json:"ssoUserId,omitempty"`   // For SSO Flow, External User ID
	ReturnTo    string             `json:"returnTo,omitempty"`    // URL to redirect after flow completion
	CreatedAt   time.Time          `json:"createdAt"`
	ExpiresAt   time.Time          `json:"expiresAt"`
}

func StoreFlow(ctx context.Context, flow FlowData) error {
	exp := time.Until(flow.ExpiresAt)
	return cached.Set(ctx, flow.ID, flow, store.WithExpiration(exp))
}

// GetFlow retrieves a flow by its ID from the cache.
//
// IMP: DO NOT RETURN nil for error if flow is not found, return a specific error instead.
func GetFlow(ctx context.Context, flowID string) (*FlowData, error) {
	if flow, err := cached.Get(ctx, flowID, new(FlowData)); err != nil {
		return nil, fmt.Errorf("failed to get flow: %w", err)
	} else {
		if f, ok := flow.(*FlowData); !ok || f == nil {
			return nil, fmt.Errorf("invalid flow data stored")
		} else {
			return f, nil
		}
	}
}
