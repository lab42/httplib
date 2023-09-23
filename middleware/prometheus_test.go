package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/lab42/httplib/internal"
	"github.com/lab42/httplib/middleware"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMiddleware(t *testing.T) {
	// Create a router with PrometheusMiddleware
	r := mux.NewRouter()
	prometheusMiddleware := middleware.NewPrometheusMiddleware(internal.DummyHandler())
	r.Use(middleware.NewPrometheusMiddleware())

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(res, req)

	// Assert the response status code is 200 OK
	assert.Equal(t, http.StatusOK, res.Code)

	// Assert Prometheus metrics were recorded as expected
	assertPrometheusMetrics(t, prometheusMiddleware)
}

func assertPrometheusMetrics(t *testing.T, prometheusMiddleware *middleware.PrometheusMiddleware) {
	// Simulate different scenarios and assert Prometheus metrics

	// Scenario 1: Successful request
	// Simulate a successful request
	// Verify that TotalRequests, ResponseStatus, and other metrics are incremented correctly
	assertMetricsIncremented(t, prometheusMiddleware, http.StatusOK)

	// Scenario 2: Request with a specific status code
	// Simulate a request with a specific status code (e.g., 404)
	// Verify that ResponseStatus metric is incremented correctly
	assertMetricsIncremented(t, prometheusMiddleware, http.StatusNotFound)

	// Scenario 4: Request with different request and response sizes
	// Simulate requests and responses with different sizes
	// Verify that RequestSize and ResponseSize metrics are observed correctly
	// Note: Adjust the sizes as needed for your test case.
	assertMetricsObserved(t, prometheusMiddleware, http.StatusOK, 1000, 2000)
	assertMetricsObserved(t, prometheusMiddleware, http.StatusOK, 500, 1000)

	// Scenario 5: Response time measurement
	// Simulate requests and measure response time
	// Verify that HttpDuration metric records response time
	assertResponseTimeMeasured(t, prometheusMiddleware)
}

func assertMetricsIncremented(t *testing.T, prometheusMiddleware *middleware.PrometheusMiddleware, statusCode int) {
	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()
	res.WriteHeader(statusCode)

	// Perform the request
	r := mux.NewRouter()
	r.Use(prometheusMiddleware.ServeHTTP)
	r.ServeHTTP(res, req)

	// Get the current route's path
	route := mux.CurrentRoute(req)
	path, _ := route.GetPathTemplate()

	// Verify that ResponseStatus metric is incremented correctly
	assert.Equal(t, 1.0, prometheusMiddleware.ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Get())

	// Verify that TotalRequests metric is incremented correctly
	assert.Equal(t, 1.0, prometheusMiddleware.TotalRequests.WithLabelValues(path).Get())

	// Verify that RequestMethods metric is incremented correctly
	assert.Equal(t, 1.0, prometheusMiddleware.RequestMethods.WithLabelValues(req.Method).Get())
}

func assertMetricsObserved(t *testing.T, prometheusMiddleware *middleware.PrometheusMiddleware, statusCode int, requestSize, responseSize int64) {
	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	req.ContentLength = requestSize
	res := httptest.NewRecorder()
	res.WriteHeader(statusCode)

	// Perform the request
	r := mux.NewRouter()
	r.Use(prometheusMiddleware.ServeHTTP)
	r.ServeHTTP(res, req)

	// Get the current route's path
	route := mux.CurrentRoute(req)
	path, _ := route.GetPathTemplate()

	// Verify that RequestSize and ResponseSize metrics are observed correctly
	assert.Equal(t, float64(requestSize), prometheusMiddleware.RequestSize.WithLabelValues(path).Get())
	assert.Equal(t, float64(responseSize), prometheusMiddleware.ResponseSize.WithLabelValues(path).Get())
}

func assertResponseTimeMeasured(t *testing.T, prometheusMiddleware *middleware.PrometheusMiddleware) {
	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Simulate a delayed handler to measure response time
	r := mux.NewRouter()
	r.Use(prometheusMiddleware.ServeHTTP)
	r.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a delay (e.g., 100 milliseconds)
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	r.ServeHTTP(res, req)

	// Get the current route's path
	route := mux.CurrentRoute(req)
	path, _ := route.GetPathTemplate()

	// Verify that HttpDuration metric records response time
	duration := prometheusMiddleware.HttpDuration.WithLabelValues(path).ObserveDuration()
	assert.True(t, duration.Seconds() >= 0.1) // Response time should be at least 100 milliseconds
}
