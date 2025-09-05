package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nbrglm/nexeres/config"
	"github.com/nbrglm/nexeres/internal/cache"
	"github.com/nbrglm/nexeres/internal/logging"
	"github.com/nbrglm/nexeres/internal/models"
	"github.com/nbrglm/nexeres/internal/tokens"
	"go.uber.org/zap"
)

const CtxAPIKeyGetter = "apiKey"
const CtxSessionToken = "sessionToken"
const CtxRefreshToken = "refreshToken"
const CtxSessionClaims = "sessionClaims"
const CtxSessionTokenClaims = "sessionTokenClaims"
const CtxSessionRefreshTokenKey = "refreshToken"
const CtxAdminToken = "adminToken"
const CtxAdminEmail = "adminEmail"

type AuthMode string

const AuthModeEither AuthMode = "either"
const AuthModeSession AuthMode = "session"
const AuthModeRefresh AuthMode = "refresh"
const AuthModeBoth AuthMode = "both"
const AuthModeAdmin AuthMode = "admin"

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
		case AuthModeAdmin:
			_, aExists := ctx.Get(CtxAdminToken)
			if !aExists {
				logging.Logger.Debug("No admin token provided")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("No admin token provided", "Provide an admin token!", http.StatusUnauthorized, nil).Filter())
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
		apiKey := strings.TrimSpace(ctx.GetHeader(tokens.NEXERES_API_KeyHeaderName))

		if apiKey == "" {
			logging.Logger.Warn("Missing API key in request")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Unauthorized access!", "Missing API key", http.StatusUnauthorized, nil).Filter())
			return
		}

		// Validate the API Key
		exists := slices.ContainsFunc(config.Security.APIKeys, func(key config.APIKeyConfig) bool {
			if apiKey == key.Key {
				ctx.Set(CtxAPIKeyGetter, key)
				return true
			}
			return false
		})

		if !exists {
			logging.Logger.Warn("Invalid API key provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("Unauthorized access!", "Invalid API key", http.StatusUnauthorized, nil).Filter())
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

		// Admin token
		adminToken := strings.TrimSpace(ctx.GetHeader(tokens.AdminTokenHeaderName))
		if adminToken != "" {
			hash := tokens.HashAdminToken(adminToken)
			sess, err := cache.GetAdminSession(ctx.Request.Context(), hash) // Just to check if it exists
			if err != nil || sess == nil {
				logging.Logger.Debug("Failed to validate admin token", zap.Error(err))
			} else {
				ctx.Set(CtxAdminToken, hash) // Store the HASH of the admin token in the context
				ctx.Set(CtxAdminEmail, sess.Email)
			}
		}

		ctx.Next()
	}
}

// ValidateSessionToken validates the provided session token and returns the claims if valid.
// It returns an error if the token is invalid or if there is an issue during validation.
func ValidateSessionToken(ctx context.Context, token string) (claims *tokens.NexeresClaims, err error) {
	// Always pass a POINTER TO THE CLAIMS STRUCT to jwt.ParseWithClaims
	// so that it can populate the claims with the parsed token data.
	// Passing a struct value will give errors like "cannot unmarshal ... into Go value of type jwt.Claims"
	// Since jwt.Claims is an interface, we need to use a pointer to a concrete type that implements it.
	c := &tokens.NexeresClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		if t.Method == jwt.SigningMethodRS256 {
			return tokens.PublicKey, nil
		} else {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithLeeway(time.Minute*5))
	if v, ok := parsedToken.Claims.(*tokens.NexeresClaims); ok && parsedToken.Valid {
		return v, nil
	}
	return nil, fmt.Errorf("invalid session token: %w", err)
}
