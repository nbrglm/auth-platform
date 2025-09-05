package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// NewAdminToken generates a new admin token and its SHA-256 hash.
// The token is a base64 URL-encoded string of 32 random bytes.
// The hash is a base64 URL-encoded SHA-256 hash of the token.
//
// Typical Approach: Return the token to the caller and store the hash securely.
// The token should be treated as a secret and only shared with authorized parties.
// The hash can be stored in a database or configuration file for later verification.
func NewAdminToken() (*string, *string, error) {
	adminTokenBytes := make([]byte, 32)
	if _, err := rand.Read(adminTokenBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to generate admin token: %w", err)
	}
	adminToken := base64.RawURLEncoding.EncodeToString(adminTokenBytes)

	adminTokenHashStr := HashAdminToken(adminToken)

	return &adminToken, &adminTokenHashStr, nil
}

// HashAdminToken hashes the provided admin token using SHA-256 and returns the base64 URL-encoded hash.
//
// Use Case: Hash the token to store OR to compare with a stored hash for verification.
func HashAdminToken(adminToken string) string {
	adminTokenHash := sha256.Sum256([]byte(adminToken))
	return base64.RawURLEncoding.EncodeToString(adminTokenHash[:])
}
