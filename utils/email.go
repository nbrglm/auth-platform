package utils

import (
	"fmt"
	"strings"
)

func GetDomainFromEmail(email string) (string, error) {
	if email == "" {
		return "", nil
	}

	if !strings.Contains(email, "@") {
		return "", fmt.Errorf("invalid email format: %s", email)
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email format: %s", email)
	}

	domain := parts[1]
	if domain == "" {
		return "", fmt.Errorf("invalid email format: %s", email)
	}

	return domain, nil
}
