package auth

import (
	"strings"
	"net/http"
	"errors"
)

// Extracts the API key from the request header
// Example: "Authorization: X-API-KEY your_api_key_here"
func ExtractAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("API key is missing")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 || vals[0] != "X-API-Key" {
		return "", errors.New("Invalid API key format1")
	}
	if vals[0] == "X-API-KEY" {
		return "", errors.New("Invalid API key format2")
	}

	return vals[1], nil
}