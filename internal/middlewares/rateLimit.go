package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/config"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

var (
	apiEndpointRateLimiter             *limiter.Limiter
	uiOpenEndpointRateLimiter          *limiter.Limiter
	uiAuthenticatedEndpointRateLimiter *limiter.Limiter
)

// InitRateLimitStore initializes the rate limit store.
// This function should be called during application startup to set up the rate limiting store.
func InitRateLimitStore() error {
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
	apiStore, err := sredis.NewStoreWithOptions(
		redisClient,
		limiter.StoreOptions{
			Prefix: "nbrglm_auth_platform_rate_limit_api",
		},
	)
	if err != nil {
		return err
	}

	uiAuthenticatedStore, err := sredis.NewStoreWithOptions(
		redisClient,
		limiter.StoreOptions{
			Prefix: "nbrglm_auth_platform_rate_limit_ui_authenticated",
		},
	)
	if err != nil {
		return err
	}

	uiOpenStore, err := sredis.NewStoreWithOptions(
		redisClient,
		limiter.StoreOptions{
			Prefix: "nbrglm_auth_platform_rate_limit_ui_open",
		},
	)
	if err != nil {
		return err
	}

	// API Endpoints
	apiRate, err := limiter.NewRateFromFormatted(config.Security.RateLimit.API.Rate)
	if err != nil {
		return err
	}
	apiEndpointRateLimiter = limiter.New(apiStore, apiRate)

	// Open UI Endpoints
	uiOpenRate, err := limiter.NewRateFromFormatted(config.Security.RateLimit.UI.OpenEndpointsRate)
	if err != nil {
		return err
	}
	uiOpenEndpointRateLimiter = limiter.New(uiOpenStore, uiOpenRate)

	// Authenticated UI Endpoints
	uiAuthenticatedRate, err := limiter.NewRateFromFormatted(config.Security.RateLimit.UI.AuthenticatedEndpointsRate)
	if err != nil {
		return err
	}
	uiAuthenticatedEndpointRateLimiter = limiter.New(uiAuthenticatedStore, uiAuthenticatedRate)
	return nil
}

// RateLimitAPIMiddleware returns a middleware that applies rate limiting
// to open API endpoints based on the configured rate limit.
// This middleware is used for endpoints that do not require authentication.
func RateLimitAPIMiddleware() gin.HandlerFunc {
	return mgin.NewMiddleware(apiEndpointRateLimiter)
}

// RateLimitUIOpenMiddleware returns a middleware that applies rate limiting
// to open UI endpoints based on the configured rate limit.
// This middleware is used for endpoints that do not require authentication.
func RateLimitUIOpenMiddleware() gin.HandlerFunc {
	return mgin.NewMiddleware(uiOpenEndpointRateLimiter)
}

// RateLimitUIAuthenticatedMiddleware returns a middleware that applies rate limiting
// to authenticated UI endpoints based on the configured rate limit.
// This middleware is used for endpoints that require user authentication.
func RateLimitUIAuthenticatedMiddleware() gin.HandlerFunc {
	return mgin.NewMiddleware(uiAuthenticatedEndpointRateLimiter)
}
