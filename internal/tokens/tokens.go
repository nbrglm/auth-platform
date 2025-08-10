// Package tokens provides functionality for managing and validating tokens used in the accounts system.
// It includes token generation, validation, and related utilities.
package tokens

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/db"
	"github.com/nbrglm/auth-platform/internal/logging"
	"github.com/nbrglm/auth-platform/opts"
	"go.uber.org/zap"
)

const RefreshTokenCookieName = "nap-refresh-tk"      // Name of the cookie used to store the refresh token.
const SessionTokenCookieName = "nap-session-tk"      // Name of the cookie used to store the session token.
const RefreshTokenHeaderName = "X-NAP-Refresh-Token" // Header name for the refresh token in API requests.
const SessionTokenHeaderName = "Authorization"       // Header name for the session token in API requests.
const RefreshTokenCookiePath = "/auth"               // Path for the refresh token cookie.
const SessionTokenCookiePath = "/"                   // Path for the session token cookie.

// RegisterHandlers registers the token-related routes with the provided Gin engine.
func RegisterHandlers(engine *gin.Engine) {
	// JWKS goes here
}

// The public/private key pair used for signing and verifying JWT tokens.
var (
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

func InitTokens() error {
	privateKeyData, err := os.ReadFile(config.JWT.PrivateKeyFile)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return errors.New("failed to parse private key PEM")
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	if privateKey == nil {
		return errors.New("private key is nil after parsing or not of type *rsa.PrivateKey")
	}

	publicKeyData, err := os.ReadFile(config.JWT.PublicKeyFile)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}
	block, _ = pem.Decode(publicKeyData)
	if block == nil {
		return errors.New("failed to parse public key PEM")
	}
	publicKeyInt, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	if pubKey, ok := publicKeyInt.(*rsa.PublicKey); ok && pubKey != nil {
		PublicKey = pubKey
	} else {
		return errors.New("public key is nil after parsing or not of type *rsa.PublicKey")
	}

	logging.Logger.Info("Tokens initialized successfully", zap.String("privateKeyFile", config.JWT.PrivateKeyFile), zap.String("publicKeyFile", config.JWT.PublicKeyFile))

	return nil
}

// RemoveTokens removes the session and refresh tokens from the Gin context and the client.
// It sets the cookies to expire immediately, effectively removing them.
// If removeSession is true, it removes the session token cookie.
// If removeRefresh is true, it removes the refresh token cookie.
// This is useful for logging out the user or clearing tokens after use.
func RemoveTokens(c *gin.Context, removeSession bool, removeRefresh bool) {
	if removeSession {
		c.SetCookie(SessionTokenCookieName, "", -1, SessionTokenCookiePath, "", !opts.Debug, true)
		logging.Logger.Debug("Removed session token cookie")
	}
	if removeRefresh {
		c.SetCookie(RefreshTokenCookieName, "", -1, RefreshTokenCookiePath, "", !opts.Debug, true)
		logging.Logger.Debug("Removed refresh token cookie")
	}
}

// SetTokens sets the session and refresh tokens in the Gin context (and in-turn, the client).
func SetTokens(c *gin.Context, tokens *Tokens) {
	c.SetSameSite(http.SameSiteStrictMode) // Before SetCookie calls

	// Empty if Debug mode, let browser decide
	refreshDomain := ""
	if !opts.Debug {
		refreshDomain = fmt.Sprintf(".%s.%s", config.Public.SubDomain, config.Public.Domain)
	}

	// Set the refresh token in a cookie
	c.SetCookie(RefreshTokenCookieName, tokens.RefreshToken, int(time.Until(tokens.RefreshTokenExpiry).Seconds()), RefreshTokenCookiePath, refreshDomain, !opts.Debug, true)

	sessionDomain := ""
	if !opts.Debug {
		sessionDomain = fmt.Sprintf(".%s", config.Public.Domain)
	}

	// Set the session ID in a cookie
	c.SetCookie(SessionTokenCookieName, tokens.SessionToken, int(time.Until(tokens.SessionTokenExpiry).Seconds()), SessionTokenCookiePath, sessionDomain, !opts.Debug, true)
}

type AuthPlatformClaims struct {
	// The registered claims from the JWT standard.
	// This struct embeds jwt.RegisteredClaims to include standard JWT claims.
	jwt.RegisteredClaims

	// Custom claims can be added here.

	// Organization information
	OrgSlug string `json:"org_slug"`
	OrgName string `json:"org_name"`
	OrgId   string `json:"org_id"`

	// User information
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	UserFname     string `json:"user_fname"`
	UserLname     string `json:"user_lname"`
	UserAvatarURL string `json:"user_avatar_url,omitempty"` // Optional user avatar URL
	UserOrgRole   string `json:"user_org_role"`
}

// Tokens represents the result of generating a new token pair.
type Tokens struct {
	// SessionId is the unique identifier for the session.
	SessionId uuid.UUID `json:"session_id"`
	// SessionToken is the generated session token.
	//
	// This is a jwt which is base64.RawURLEncoding encoded.
	// YOU NEED TO DECODE IT WHILE RETRIEVING IT FROM THE COOKIES/CLIENT.
	// DO NOT USE IT AS IS. VALIDATION WILL FAIL WITHOUT DECODING.
	// ONLY WHEN DECODED, YOU SHOULD PASS IT TO THE THINGS THAT REQUIRE THE SESSION TOKEN.
	SessionToken string `json:"session_token"`
	// SessionTokenExpiry is the expiration time of the session token.
	SessionTokenExpiry time.Time `json:"session_token_expiry"`
	// RefreshToken is the generated refresh token.
	//
	// This is base64.RawURLEncoding encoded.
	// DO NOT DECODE IT WHILE RETRIEVING IT FROM THE COOKIES/CLIENT.
	RefreshToken string `json:"refresh_token"`
	// RefreshTokenExpiry is the expiration time of the refresh token.
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}

// GenerateTokens generates a new session and refresh token pair for the given user ID and claims.
//
// NOTE: This function will NOT store the tokens in the database.
func GenerateTokens(userId uuid.UUID, claims AuthPlatformClaims) (*Tokens, error) {
	now := time.Now().UTC()

	sessionId, err := uuid.NewV7()
	if err != nil {
		logging.Logger.Error("Failed to generate session ID", zap.Error(err))
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Calculate expiration durations
	sessionTokenExpirationDuration := time.Duration(config.JWT.SessionTokenExpiration) * time.Second
	refreshTokenExpDuration := time.Duration(config.JWT.RefreshTokenExpiration) * time.Second

	// Set the standard claims as per RFC 7519 JWT specification.
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    config.Public.GetBaseURL(),
		Subject:   userId.String(),
		Audience:  jwt.ClaimStrings(config.JWT.Audiences),
		ExpiresAt: jwt.NewNumericDate(now.Add(sessionTokenExpirationDuration)),
		NotBefore: jwt.NewNumericDate(now.Add(-(time.Minute * 5))), // 5 minutes before the token is valid
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        sessionId.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	sessionToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign session token: %w", err)
	}
	sessionToken = base64.RawURLEncoding.EncodeToString([]byte(sessionToken)) // Encode to base64 URL-safe string

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(refreshTokenBytes)

	return &Tokens{
		SessionId:          sessionId,
		SessionToken:       sessionToken,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: now.Add(refreshTokenExpDuration),
		SessionTokenExpiry: now.Add(sessionTokenExpirationDuration),
	}, nil
}

// RefreshSession refreshes the session by generating a new session and refresh token pair.
// All the non-standard claims have to be set before passing in the claims parameter.
// It takes the old session and claims as the parameter, and gives new tokens.
func RefreshSessionTokens(session db.Session, claims AuthPlatformClaims) (*Tokens, error) {
	now := time.Now().UTC()

	// Calculate expiration durations
	sessionTokenExpiryDuration := time.Duration(config.JWT.SessionTokenExpiration) * time.Second
	refreshTokenExpDuration := time.Duration(config.JWT.RefreshTokenExpiration) * time.Second

	// Set the standard claims as per RFC 7519 JWT specification.
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    config.Public.GetBaseURL(),
		Subject:   session.UserID.String(),
		Audience:  jwt.ClaimStrings(config.JWT.Audiences),
		ExpiresAt: jwt.NewNumericDate(now.Add(sessionTokenExpiryDuration)),
		NotBefore: jwt.NewNumericDate(now.Add(-(time.Minute * 5))), // 5 minutes before the token is valid
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        session.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	sessionToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign session token: %w", err)
	}

	sessionToken = base64.RawURLEncoding.EncodeToString([]byte(sessionToken)) // Encode to base64 URL-safe string

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(refreshTokenBytes)

	return &Tokens{
		SessionId:          session.ID,
		SessionToken:       sessionToken,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: now.Add(refreshTokenExpDuration),
		SessionTokenExpiry: now.Add(sessionTokenExpiryDuration),
	}, nil
}

// HashTokens hashes the session and refresh tokens using SHA-256.
// It returns the hashed values as strings (And empty strings if the input tokens are empty).
func HashTokens(tokens *Tokens) (string, string) {
	if tokens == nil {
		return "", ""
	}
	sessionHash := sha256.Sum256([]byte(tokens.SessionToken))
	refreshHash := sha256.Sum256([]byte(tokens.RefreshToken))
	return hex.EncodeToString(sessionHash[:]), hex.EncodeToString(refreshHash[:])
}

func HasTokenBeenRevoked(ctx context.Context, q *db.Queries, sessionId uuid.UUID) (bool, error) {
	_, err := q.GetSessionByID(ctx, sessionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil // Session not found, meaning it has been revoked
		}
		return false, fmt.Errorf("failed to get session: %w", err)
	}
	return false, nil
}
