package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/internal/logging"
	"github.com/nbrglm/auth-platform/internal/models"
	"github.com/nbrglm/auth-platform/internal/tokens"
	"go.uber.org/zap"
)

const CtxAPIKeyGetter = "apiKey"
const CtxSessionToken = "sessionToken"
const CtxRefreshToken = "refreshToken"
const CtxSessionClaims = "sessionClaims"
const CtxSessionTokenClaims = "sessionTokenClaims"
const CtxSessionRefreshTokenKey = "refreshToken"

type AuthMode string

const AuthModeEither AuthMode = "either"
const AuthModeSession AuthMode = "session"
const AuthModeRefresh AuthMode = "refresh"
const AuthModeBoth AuthMode = "both"

func RequireAuth(mode AuthMode) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if at least one of the tokens is present based on the mode
		switch mode {
		case AuthModeEither:
			_, sExists := ctx.Get(CtxSessionToken)
			_, rExists := ctx.Get(CtxRefreshToken)
			if !sExists && !rExists {
				logging.Logger.Debug("No session or refresh token provided")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("No session or refresh token provided", "Provide either a session token or a refresh token", http.StatusUnauthorized, nil).Filter())
				return
			}
		case AuthModeSession:
			_, sExists := ctx.Get(CtxSessionToken)
			if !sExists {
				logging.Logger.Debug("No session token provided")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("No session token provided", "Provide a session token!", http.StatusUnauthorized, nil).Filter())
				return
			}
		case AuthModeRefresh:
			_, rExists := ctx.Get(CtxRefreshToken)
			if !rExists {
				logging.Logger.Debug("No refresh token provided")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("No refresh token provided", "Provide a refresh token!", http.StatusUnauthorized, nil).Filter())
				return
			}
		case AuthModeBoth:
			_, sExists := ctx.Get(CtxSessionToken)
			_, rExists := ctx.Get(CtxRefreshToken)
			if !sExists || !rExists {
				logging.Logger.Debug("No session or refresh token provided")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("No session or refresh token provided", "Provide both a session token and a refresh token!", http.StatusUnauthorized, nil).Filter())
				return
			}
		default:
			logging.Logger.Error("Invalid auth mode provided to RequireAuth middleware", zap.String("mode", string(mode)))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("Internal server error", "An internal server error occurred", http.StatusInternalServerError, nil).Filter())
			return
		}
		ctx.Next()
	}
}

func APIKeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := strings.TrimSpace(ctx.GetHeader(tokens.NAP_API_KeyHeaderName))

		if apiKey == "" {
			logging.Logger.Warn("Missing API key in request")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			return
		}

		// Validate the API Key
		for _, key := range config.Security.APIKeys {
			if apiKey == key.Key {
				ctx.Set(CtxAPIKeyGetter, key)
				break
			}
		}
		if _, exists := ctx.Get(CtxAPIKeyGetter); !exists {
			logging.Logger.Warn("Invalid API key provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		// Session token
		sessionToken := strings.TrimSpace(ctx.GetHeader(tokens.SessionTokenHeaderName))
		if sessionToken != "" {
			bytes, err := base64.RawURLEncoding.DecodeString(sessionToken)
			if err != nil {
				logging.Logger.Debug("Failed to decode session token", zap.Error(err))
			} else {
				// Validate the session token
				claims, err := ValidateSessionToken(ctx.Request.Context(), string(bytes))
				if err != nil {
					logging.Logger.Debug("Failed to validate session token", zap.Error(err))
				} else if claims != nil {
					ctx.Set(CtxSessionToken, sessionToken)
					ctx.Set(CtxSessionTokenClaims, claims)
				}
			}
		}

		// Refresh token
		refreshToken := strings.TrimSpace(ctx.GetHeader(tokens.RefreshTokenHeaderName))
		if refreshToken != "" {
			ctx.Set(CtxRefreshToken, refreshToken)
		}
	}
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
