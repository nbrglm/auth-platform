package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nbrglm/auth-platform/internal/logging"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"go.uber.org/zap"
)

const CtxSessionExistsKey = "sessionExists"
const CtxSessionRefreshKey = "sessionRefresh"
const CtxSessionTokenKey = "sessionToken"
const CtxSessionTokenClaimsKey = "sessionTokenClaims"
const CtxSessionRefreshTokenKey = "refreshToken"

func APIKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Put the token info, like claims, the tokens themselves, etc. into the context.
		session, refresh := ExpandContextWithInfo(ctx)

		ctx.Set(CtxSessionExistsKey, session)
		ctx.Set(CtxSessionRefreshKey, refresh)

		// Continue to the next middleware/handler
		ctx.Next()
	}
}

// ExpandContextWithInfo extracts auth tokens from request headers or cookies and adds them to Gin context.
// - Validates session tokens from "Authorization" header or "nap-session-tk" cookie
// - Stores refresh tokens from "X-NAP-Refresh-Token" header or "nap-refresh-tk" cookie (without validation)
// Returns:
// - session: true if valid session token found
// - refresh: true if refresh token found
func ExpandContextWithInfo(c *gin.Context) (session bool, refresh bool) {
	getToken := func(headerName, cookieName string) string {
		// Check header first
		token := strings.TrimSpace(c.Request.Header.Get(headerName))
		if token != "" {
			logging.Logger.Debug("Token found in header",
				zap.String("header", headerName))
			return token
		}

		// If not in header, check cookie
		logging.Logger.Debug("Token not found in header, checking cookie",
			zap.String("header", headerName),
			zap.String("cookie", cookieName))

		cookieVal, err := c.Cookie(cookieName)
		if err != nil {
			logging.Logger.Debug("Cookie not found",
				zap.String("cookie", cookieName),
				zap.Error(err))
			return ""
		}
		return strings.TrimSpace(cookieVal)
	}

	sessionToken := getToken(tokens.SessionTokenHeaderName, tokens.SessionTokenCookieName)
	if sessionToken != "" {
		bytes, err := base64.RawURLEncoding.DecodeString(sessionToken)
		if err != nil {
			logging.Logger.Debug("Failed to decode session token", zap.Error(err))
			return false, false
		}
		// Validate the session token
		claims, err := ValidateSessionToken(c.Request.Context(), string(bytes))
		if err != nil {
			logging.Logger.Debug("Failed to validate session token", zap.Error(err))
		}

		if claims != nil {
			c.Set(CtxSessionTokenKey, sessionToken)
			c.Set(CtxSessionTokenClaimsKey, claims)
			session = true
		}
	}

	// Refresh token handling
	refreshToken := getToken(tokens.RefreshTokenHeaderName, tokens.RefreshTokenCookieName)
	if refreshToken != "" {
		c.Set(CtxSessionRefreshTokenKey, refreshToken)
		refresh = true
	}

	return
}

// ValidateSessionToken validates the provided session token and returns the claims if valid.
// It returns an error if the token is invalid or if there is an issue during validation.
func ValidateSessionToken(ctx context.Context, token string) (claims *tokens.AuthPlatformClaims, err error) {
	// Always pass a POINTER TO THE CLAIMS STRUCT to jwt.ParseWithClaims
	// so that it can populate the claims with the parsed token data.
	// Passing a struct value will give errors like "cannot unmarshal ... into Go value of type jwt.Claims"
	// Since jwt.Claims is an interface, we need to use a pointer to a concrete type that implements it.
	c := &tokens.AuthPlatformClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		if t.Method == jwt.SigningMethodRS256 {
			return tokens.PublicKey, nil
		} else {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithLeeway(time.Minute*5))
	if v, ok := parsedToken.Claims.(*tokens.AuthPlatformClaims); ok && parsedToken.Valid {
		return v, nil
	}
	return nil, fmt.Errorf("invalid session token: %w", err)
}
