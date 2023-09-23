package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/lab42/httplib/internal"
	"github.com/lab42/httplib/middleware"
)

func TestRateLimitMiddleware_Redis(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	// Create a RateLimitMiddleware instance for testing.
	middleware := middleware.RateLimitMiddleware{
		Next:     http.HandlerFunc(internal.DummyHandler),
		Every:    2,
		Redis:    rdb,
		InMemory: nil,
	}

	// Create a test request.
	req := httptest.NewRequest("GET", "/", nil)

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Define the IP address for testing.
	ip := "127.0.0.1"

	// Create a Redis transaction pipeline that increments tokens and sets the last refill time.
	mock.ExpectHIncrBy("ratelimit:"+ip, "tokens", int64(2)).SetVal(2)

	mock.ExpectHGet("ratelimit:"+ip, "lastRefill").SetVal(
		strconv.FormatInt(time.Now().Add(-5*time.Second).Unix(), 10),
	)
	mock.ExpectHSet("ratelimit:"+ip, "lastRefill").SetVal(1)
	mock.ExpectHIncrBy("ratelimit:"+ip, "tokens", int64(2)).SetVal(3)

	// Execute the middleware 3 times in quick succession.
	for i := 0; i < 3; i++ {
		middleware.ServeHTTP(rr, req)
	}

	// Check if the response status code is 200 for the first two requests.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", status)
	}

	// Check if the response status code is 429 (Too Many Requests) for the third request.
	// Also, check if the response body contains the rate limit exceeded message.
	if status := rr.Code; status != http.StatusTooManyRequests {
		t.Errorf("Expected status code 429 (Too Many Requests), but got %d", status)
	}

	expectedError := "Rate limit exceeded"
	if rr.Body.String() != expectedError {
		t.Errorf("Expected response body '%s', but got '%s'", expectedError, rr.Body.String())
	}
}

func TestRateLimitMiddleware_InMemory(t *testing.T) {
	// Create a RateLimitMiddleware instance with in-memory storage for testing.
	middleware := middleware.RateLimitMiddleware{
		Next:     http.HandlerFunc(dummyHandler), // Define a dummy handler for testing.
		Every:    2,                              // Allow 2 requests per second.
		Redis:    nil,
		InMemory: &sync.Map{},
	}

	// Create a test request.
	req := httptest.NewRequest("GET", "/", nil)

	// Create a response recorder to capture the responses.
	rr := httptest.NewRecorder()

	// Execute the middleware 3 times in quick succession.
	for i := 0; i < 3; i++ {
		middleware.ServeHTTP(rr, req)
	}

	// Check if the response code for the first two requests is OK.
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// The third request should be rate-limited and return a 429 status code.
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, but got %d", http.StatusTooManyRequests, rr.Code)
	}
}

func TestRateLimitMiddleware_RateLimitExceededError(t *testing.T) {
	// Create a mock Redis client that returns an error for testing rate limit exceeded scenario.
	rdb, mock := redismock.NewClientMock()

	// Create a RateLimitMiddleware instance for testing.
	middleware := middleware.RateLimitMiddleware{
		Next:     http.HandlerFunc(dummyHandler), // Define a dummy handler for testing.
		Every:    2,                              // Allow 2 requests per second.
		Redis:    rdb,
		InMemory: &sync.Map{},
	}

	// Create a test request.
	req := httptest.NewRequest("GET", "/", nil)

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Create a mock Redis transaction pipeline that returns an error.
	mock.ExpectHIncrBy("ratelimit:127.0.0.1", "tokens", int64(2)).SetErr(errors.New("mock error"))
	mock.ExpectHGet("ratelimit:127.0.0.1", "lastRefill").SetVal("1633290000")

	// Execute the middleware 3 times in quick succession.
	for i := 0; i < 3; i++ {
		middleware.ServeHTTP(rr, req)
	}

	// Check if the response body contains the rate limit exceeded error message.
	expectedError := "Rate limit exceeded"
	if strings.Trim(rr.Body.String(), "\n") != expectedError {
		t.Errorf("Expected response body '%s', but got '%s'", expectedError, rr.Body.String())
	}
}
