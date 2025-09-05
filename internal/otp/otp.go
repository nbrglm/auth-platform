package otp

import (
	"crypto/rand"
	"fmt"
)

// NewAlphaNumericOTP generates an Alphanumeric OTP of the specified length.
//
// Use Case: This function can be used to generate OTPs for two-factor authentication (2FA),
// password resets, or any other scenario where a short-lived code is required for verification.
func NewAlphaNumericOTP(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	// size of charset (62)
	charsetLen := byte(len(charset))
	// max random byte we can safely accept without bias
	// (largest multiple of charsetLen < 256)
	maxRandom := byte(256 - (256 % int(charsetLen)))

	for i := range length {
		var b [1]byte
		for {
			// read one random byte
			if _, err := rand.Read(b[:]); err != nil {
				return "", fmt.Errorf("failed to generate OTP: %w", err)
			}
			// only accept if it's in unbiased range
			if b[0] < maxRandom {
				result[i] = charset[b[0]%charsetLen]
				break
			}
			// otherwise loop again
		}
	}

	return string(result), nil
}
