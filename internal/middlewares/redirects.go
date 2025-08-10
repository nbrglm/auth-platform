package middlewares

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewRedirectIfAuthenticatedMiddleware creates a middleware that redirects authenticated users away from login pages.
//
// Use this middleware JUST before the route, and not with the `Use` method.
func NewRedirectIfAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, og, safe := GetRedirectURL(c)

		// Check if the user is authenticated
		if c.GetBool(CtxSessionExistsKey) {
			// If authenticated, redirect to the home page
			c.Redirect(http.StatusSeeOther, u)
			c.Abort() // Stop further processing
			return
		} else if c.GetBool(CtxSessionRefreshKey) {
			if safe {
				c.Redirect(http.StatusSeeOther, "/auth/s/refresh?returnTo="+url.QueryEscape(og))
			} else {
				c.Redirect(http.StatusSeeOther, "/auth/s/refresh")
			}
			c.Abort()
			return
		}
		c.Next() // Continue to the next middleware or handler
	}
	// This middleware can be used to redirect authenticated users away from login pages.
}

// GetRedirectURLOriginalOrFallback redirects the user to the original redirect URL if it is safe.
// If the original redirect URL is not safe, it redirects to the default unsafe-redirect page.
// If the original redirect URL is empty, it defaults to "/unsafe-redirect".
// It returns the URL redirected to.
func GetRedirectURLOriginalOrFallback(c *gin.Context) string {
	u, og, safe := GetRedirectURL(c)

	// If the redirect URL is safe, redirect to it directly.
	if safe {
		return u
	}

	// If the original redirect URL is provided, redirect to unsafe-redirect page with the original URL.
	if og != "" {
		return "/unsafe-redirect?original=" + url.QueryEscape(og)
	}

	// If no valid redirect URL is found, redirect to the default unsafe-redirect page.
	return "/unsafe-redirect"
}

// GetRedirectURLWithReturnTo redirects the user to a specified path with a returnTo query parameter.
// If the redirect URL is safe, it appends the original redirect URL as a query parameter.
// If the original redirect URL is empty, it defaults to "/unsafe-redirect".
// It returns the URL it redirected to.
func GetRedirectURLWithReturnTo(c *gin.Context, path string) string {
	u, og, safe := GetRedirectURL(c)

	redirectTo := path

	if safe {
		redirectTo += "?returnTo=" + url.QueryEscape(u)
	} else if og != "" {
		redirectTo += "?returnTo=" + url.QueryEscape(og)
	}

	return redirectTo
}

// GetRedirectURL retrieves the redirect URL from the context and formats it based on whether it is safe or unsafe.
// It returns the formatted redirect URL and the original redirect URL.
// If the original redirect URL is empty, it defaults to "/unsafe-redirect".
// If the unsafe flag is set, it appends the original redirect URL as a query parameter
func GetRedirectURL(c *gin.Context) (string, string, bool) {
	originalRedirectURL := strings.TrimSpace(c.GetString(ReturnToURLKey))
	unsafe := c.GetBool(ReturnToURLUnsafeKey)
	newRedirectURL := "/unsafe-redirect"

	if originalRedirectURL != "" {
		if unsafe {
			newRedirectURL += "?original=" + url.QueryEscape(originalRedirectURL)
		} else {
			newRedirectURL = originalRedirectURL
		}
	}

	return newRedirectURL, originalRedirectURL, !unsafe
}
