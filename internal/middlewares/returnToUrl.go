package middlewares

import (
	"net/url"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/internal/logging"
	"go.uber.org/zap"
)

const ReturnToURLKey = "returnTo"
const ReturnToURLUnsafeKey = "unsafeReturnTo"

// ReturnToURLMiddleware is a middleware that checks for a "returnTo" query parameter in the request.
// If present, it validates the URL and stores it in the context.
// if present, but not under the allowed domains, it will set it in the context as follows:
// unsafe://<returnTo URL>
//
// This middleware should be used before any handler that needs to access the returnTo URL.
// It is typically used in routes that require a return URL, such as login or signup.
func ReturnToURLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		returnTo := c.Query("returnTo")
		if returnTo == "" {
			c.Next()
			return
		}

		u, err := url.Parse(returnTo)
		if err != nil {
			// Skip if the URL is invalid
			c.Next()
			return
		}

		if !u.IsAbs() {
			path := strings.TrimSpace(u.Path)
			if strings.HasPrefix(path, "/") {
				// relative path case (safe, since it stays on same host)
				c.Set(ReturnToURLKey, path)
				c.Set(ReturnToURLUnsafeKey, false)
				logging.Logger.Debug("relative returnTo URL detected", zap.String("returnTo", path))
			}
			// if not a valid relative path, skip
			c.Next()
			return
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			// Skip if the URL scheme is not http or https
			c.Next()
			return
		}

		// validate the returnTo URL
		if slices.ContainsFunc(config.Public.Redirects.Domains, func(domain string) bool {
			return strings.HasSuffix(u.Hostname(), domain)
		}) {
			// If the returnTo URL is valid, set it in the context
			c.Set(ReturnToURLKey, returnTo)
			c.Set(ReturnToURLUnsafeKey, false)
		} else {
			c.Set(ReturnToURLKey, returnTo)
			c.Set(ReturnToURLUnsafeKey, true)
		}
		// Continue processing the request
		c.Next()
	}
}
