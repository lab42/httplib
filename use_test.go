package httplib_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lab42/httplib"
)

func TestUseMiddleware(t *testing.T) {
	// Create a new ServeMux and a simple handler for testing.
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Define a middleware function for testing.
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-Header", "MiddlewareApplied")
			next.ServeHTTP(w, r)
		})
	}

	// Apply the middleware using the Use function.
	handler := httplib.Use(mux, middleware)

	// Create a test request.
	req := httptest.NewRequest("GET", "/test", nil)

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Send the request through the middleware chain.
	handler.ServeHTTP(rr, req)

	// Check if the middleware header was set.
	if rr.Header().Get("X-Middleware-Header") != "MiddlewareApplied" {
		t.Errorf("Expected X-Middleware-Header to be 'MiddlewareApplied', but got '%s'", rr.Header().Get("X-Middleware-Header"))
	}

	// Check the response body.
	expectedResponse := "Hello, World!"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response '%s', but got '%s'", expectedResponse, rr.Body.String())
	}
}
