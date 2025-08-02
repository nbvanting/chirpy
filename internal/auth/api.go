package auth

import (
	"fmt"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	const prefix = "ApiKey "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid authorization header format")
	}
	apiKey := authHeader[len(prefix):]
	if apiKey == "" {
		return "", fmt.Errorf("API key missing")
	}
	return apiKey, nil
}