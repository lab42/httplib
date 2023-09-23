package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lab42/httplib/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddlewareWithPanic(t *testing.T) {
	// Create a test handler that panics.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Initialize the RecoveryMiddleware with the test handler.
	middleware := middleware.NewRecoveryMiddleware(testHandler)

	// Create a test request.
	req := httptest.NewRequest("GET", "http://example.com/test", nil)

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Serve the request using the middleware.
	middleware.ServeHTTP(rr, req)

	// Ensure that the response status code is 500 (Internal Server Error).
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRecoveryMiddlewareWithoutPanic(t *testing.T) {
	// Create a test handler that does not panic.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Initialize the RecoveryMiddleware with the test handler.
	middleware := middleware.NewRecoveryMiddleware(testHandler)

	// Create a test request.
	req := httptest.NewRequest("GET", "http://example.com/test", nil)

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Serve the request using the middleware.
	middleware.ServeHTTP(rr, req)

	// Ensure that the response status code is 200 (OK).
	assert.Equal(t, http.StatusOK, rr.Code)
}
