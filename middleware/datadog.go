package middleware

import (
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// DatadogMiddleware is a middleware that sends request metrics to Datadog.
type DatadogMiddleware struct {
	Next             http.Handler
	StatsDClient     *statsd.Client
	StatsDSampleRate float64 // Sample rate for collecting metrics (0.0 - 1.0).
}

// NewDatadogMiddleware creates a new DatadogMiddleware instance.
func NewDatadogMiddleware(next http.Handler, statsdClient *statsd.Client, sampleRate float64) *DatadogMiddleware {
	return &DatadogMiddleware{
		Next:             next,
		StatsDClient:     statsdClient,
		StatsDSampleRate: sampleRate,
	}
}

// ServeHTTP is the middleware handler function that collects and sends request metrics to Datadog.
func (m *DatadogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Call the next handler in the chain.
	m.Next.ServeHTTP(w, r)

	// Calculate request duration.
	duration := time.Since(startTime)

	// Send commonly used request metrics to Datadog.
	m.StatsDClient.Histogram("http.request.duration", float64(duration.Seconds()), []string{}, m.StatsDSampleRate)
	m.StatsDClient.Incr("http.request.count", []string{}, m.StatsDSampleRate)
	m.StatsDClient.Incr("http.response.status."+http.StatusText(w.(*datadogResponseWriter).statusCode), []string{}, m.StatsDSampleRate)

	// Log the request method as a tag.
	m.StatsDClient.Incr("http.request.method."+r.Method, []string{}, m.StatsDSampleRate)

	// Log the request path as a tag.
	m.StatsDClient.Incr("http.request.path."+r.URL.Path, []string{}, m.StatsDSampleRate)
}

// responseWriter is a custom http.ResponseWriter that tracks response status code.
type datadogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new responseWriter instance.
func NewDatadogResponseWriter(w http.ResponseWriter) *datadogResponseWriter {
	return &datadogResponseWriter{w, http.StatusOK}
}

// WriteHeader intercepts the WriteHeader method to track response status code.
func (rw *datadogResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
