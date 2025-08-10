package utils

import (
	"encoding/base64"
	"fmt"
)

func DecodeB64Key(encodedKey string) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key: %w", err)
	}
	return decodedKey, nil
}

func EncodeB64Key(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}
