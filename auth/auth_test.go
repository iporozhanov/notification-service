package auth

import (
	"net/http"
	"testing"
)

func TestHTTPAPIKeyAuth_AuthenticateRequest(t *testing.T) {
	// Create a new instance of HTTPAPIKeyAuth
	auth := NewHTTPAPIKeyAuth("mySecretKey")

	// Create a new request with the AUTH_TOKEN header set to the secret key
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("AUTH_TOKEN", "mySecretKey")

	// Call the AuthenticateRequest method and check the result
	_, err := auth.AuthenticateRequest(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Create a new request with an invalid AUTH_TOKEN header
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("AUTH_TOKEN", "invalidKey")

	// Call the AuthenticateRequest method and check the result
	_, err = auth.AuthenticateRequest(req)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
