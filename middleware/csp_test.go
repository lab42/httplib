package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lab42/httplib/middleware"
)

func TestCSPMiddleware(t *testing.T) {
	// Create a sample HTTP handler for testing.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create a test request.
	req := httptest.NewRequest("GET", "/", nil)

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Define a CSP header value for testing.
	cspHeaderValue := "default-src 'self'; script-src 'unsafe-inline'"

	// Create a CSPMiddleware instance.
	middleware := middleware.NewCSPMiddleware(handler, cspHeaderValue)

	// Execute the middleware.
	middleware.ServeHTTP(rr, req)

	// Check if the CSP header was set correctly in the response.
	if rr.Header().Get("Content-Security-Policy") != cspHeaderValue {
		t.Errorf("Expected CSP header '%s', but got '%s'", cspHeaderValue, rr.Header().Get("Content-Security-Policy"))
	}

	// Check the response body.
	expectedResponse := "Hello, World!"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response '%s', but got '%s'", expectedResponse, rr.Body.String())
	}
}
