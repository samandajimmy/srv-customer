package nhttp

import (
	"encoding/base64"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"strings"
)

func ExtractAuthValue(prefix string, str string) (string, error) {
	if str == "" {
		return "", ncore.NewError("authorization value is empty")
	}

	// Extract token
	tokens := strings.Split(str, " ")
	if len(tokens) != 2 {
		return "", ncore.NewError("token is malformed")
	}

	// Check prefix
	if tokens[0] != prefix {
		return "", ncore.NewError("unexpected prefix")
	}

	return tokens[1], nil
}

func ExtractBasicAuth(r *http.Request) (username string, password string, err error) {
	// Get header
	authHeader := r.Header.Get(AuthorizationHeader)

	// Extract base64 encoded
	encAuth, err := ExtractAuthValue("Basic", authHeader)
	if err != nil {
		return "", "", err
	}

	// Decode base64 auth
	decAuth, err := base64.StdEncoding.DecodeString(encAuth)
	if err != nil {
		return "", "", err
	}

	// Split row by : delimiter
	tokens := strings.Split(string(decAuth), ":")
	if len(tokens) != 2 {
		return "", "", ncore.NewError("malformed decoded Basic Auth")
	}

	return tokens[0], tokens[1], nil
}

func ExtractBearerAuth(r *http.Request) (token string, err error) {
	// Get header
	authHeader := r.Header.Get(AuthorizationHeader)

	// Extract base64 encoded
	bearerToken, err := ExtractAuthValue("Bearer", authHeader)
	if err != nil {
		return "", err
	}

	return bearerToken, nil
}
