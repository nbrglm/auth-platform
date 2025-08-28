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
	rateLimiter *limiter.Limiter
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
	rateLimitStore, err := sredis.NewStoreWithOptions(
		redisClient,
		limiter.StoreOptions{
			Prefix: "nbrglm_auth_platform_rate_limit",
		},
	)
	if err != nil {
		return err
	}

	rate, err := limiter.NewRateFromFormatted(config.Security.RateLimit.Rate)
	if err != nil {
		return err
	}

	rateLimiter = limiter.New(rateLimitStore, rate)
	return nil
}

// RateLimitMiddleware returns a middleware that applies rate limiting
// to open API endpoints based on the configured rate limit.
// This middleware is used for endpoints that do not require authentication.
func RateLimitMiddleware() gin.HandlerFunc {
	return mgin.NewMiddleware(rateLimiter)
}
