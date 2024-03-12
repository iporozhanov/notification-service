package auth

import (
	"fmt"
	"net/http"
)

// Auth represents an interface for request authentication.
type HTTPAPIKeyAuth struct {
	secret string
}

func NewHTTPAPIKeyAuth(secret string) *HTTPAPIKeyAuth {
	return &HTTPAPIKeyAuth{
		secret: secret,
	}
}

// AuthenticateRequest authenticates an HTTP request.
func (a *HTTPAPIKeyAuth) AuthenticateRequest(r *http.Request) (string, error) {
	tokenStr := r.Header.Get("AUTH_TOKEN")
	if tokenStr == "" {
		return "", fmt.Errorf("no token provided")
	}

	if tokenStr != a.secret {
		return "", fmt.Errorf("invalid token")
	}

	return "", nil
}
