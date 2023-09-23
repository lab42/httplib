package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lab42/httplib/middleware"
)

func TestCSRFMiddleware(t *testing.T) {
	// Create a sample HTTP handler for testing.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create a test request.
	req := httptest.NewRequest("GET", "/", nil)

	// Set a CSRF token in the request header for testing.
	csrfToken := "abc123"
	req.Header.Set("X-CSRF-Token", csrfToken)

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Create a CSRFMiddleware instance.
	csrfMiddleware := middleware.NewCSRFMiddleware(handler, "X-CSRF-Token", "csrfCookie", "csrfParam", csrfToken, http.StatusForbidden)

	// Execute the middleware.
	csrfMiddleware.ServeHTTP(rr, req)

	// Check if the response status code is OK (should pass CSRF validation).
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body.
	expectedResponse := "Hello, World!"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response '%s', but got '%s'", expectedResponse, rr.Body.String())
	}

	// Test case: Invalid CSRF token, should return Forbidden status.
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-CSRF-Token", "invalidtoken")
	rr = httptest.NewRecorder()

	csrfMiddleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, but got %d", http.StatusForbidden, rr.Code)
	}
}
