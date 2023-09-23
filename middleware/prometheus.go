package middleware

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMiddleware is a middleware that collects Prometheus metrics for HTTP requests.
type PrometheusMiddleware struct {
	Next                        http.Handler
	TotalRequests               *prometheus.CounterVec
	ResponseStatus              *prometheus.CounterVec
	HttpDuration                *prometheus.HistogramVec
	RequestMethods              *prometheus.CounterVec
	RequestSize                 *prometheus.HistogramVec
	ResponseSize                *prometheus.HistogramVec
	HttpDurationByMethod        *prometheus.HistogramVec
	HttpResponseTimePercentiles *prometheus.SummaryVec
}

// NewPrometheusMiddleware creates a new PrometheusMiddleware instance.
func NewPrometheusMiddleware(next http.Handler) *PrometheusMiddleware {
	// Create Prometheus metrics
	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of requests",
		},
		[]string{"path"},
	)

	responseStatus := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests",
	}, []string{"path"})

	requestMethods := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_methods_total",
			Help: "Number of requests by HTTP method",
		},
		[]string{"method"},
	)

	requestSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_size_bytes",
			Help: "Size of incoming HTTP requests",
		},
		[]string{"path"},
	)

	responseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_size_bytes",
			Help: "Size of HTTP responses",
		},
		[]string{"path"},
	)

	httpDurationByMethod := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds_by_method",
		Help: "Duration of HTTP requests by method",
	}, []string{"method"})

	httpResponseTimePercentiles := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_response_time_percentiles",
		Help:       "Response time percentiles",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	}, []string{"path"})

	// Create PrometheusMiddleware instance
	prometheusMiddleware := &PrometheusMiddleware{
		Next:                        next,
		TotalRequests:               totalRequests,
		ResponseStatus:              responseStatus,
		HttpDuration:                httpDuration,
		RequestMethods:              requestMethods,
		RequestSize:                 requestSize,
		ResponseSize:                responseSize,
		HttpDurationByMethod:        httpDurationByMethod,
		HttpResponseTimePercentiles: httpResponseTimePercentiles,
	}

	// Register Prometheus metrics
	prometheus.MustRegister(
		prometheusMiddleware.TotalRequests,
		prometheusMiddleware.ResponseStatus,
		prometheusMiddleware.HttpDuration,
		prometheusMiddleware.RequestMethods,
		prometheusMiddleware.RequestSize,
		prometheusMiddleware.ResponseSize,
		prometheusMiddleware.HttpDurationByMethod,
		prometheusMiddleware.HttpResponseTimePercentiles,
	)

	return prometheusMiddleware
}

// ServeHTTP is the middleware handler function that collects Prometheus metrics.
func (m *PrometheusMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the current route's path
	route := mux.CurrentRoute(r)
	path, _ := route.GetPathTemplate()

	// Start measuring response time
	timer := prometheus.NewTimer(m.HttpDuration.WithLabelValues(path))

	// Create a custom response writer to track response size
	customRW := NewPrometheusResponseWriter(w)

	// Call the next handler in the chain
	m.Next.ServeHTTP(customRW, r)

	// Get the response status code
	statusCode := customRW.statusCode

	// Increment response status counter
	m.ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()

	// Increment total requests counter for the specific path
	m.TotalRequests.WithLabelValues(path).Inc()

	// Increment request methods counter for the HTTP method used
	m.RequestMethods.WithLabelValues(r.Method).Inc()

	// Calculate request size
	requestSize := r.ContentLength
	if requestSize < 0 {
		requestSize = 0
	}
	m.RequestSize.WithLabelValues(path).Observe(float64(requestSize))

	// Calculate response size
	responseSize := customRW.Size()
	m.ResponseSize.WithLabelValues(path).Observe(float64(responseSize))

	// Calculate response time by HTTP method
	m.HttpDurationByMethod.WithLabelValues(r.Method).Observe(timer.ObserveDuration().Seconds())

	// Calculate response time percentiles
	m.HttpResponseTimePercentiles.WithLabelValues(path).Observe(timer.ObserveDuration().Seconds())
}

// responseWriter is a custom http.ResponseWriter that tracks response status code.
type prometheusResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

// NewResponseWriter creates a new responseWriter instance.
func NewPrometheusResponseWriter(w http.ResponseWriter) *prometheusResponseWriter {
	return &prometheusResponseWriter{w, http.StatusOK, 0}
}

// WriteHeader intercepts the WriteHeader method to track response status code.
func (rw *prometheusResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write intercepts the Write method to track response size.
func (rw *prometheusResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

// Size returns the accumulated response size.
func (rw *prometheusResponseWriter) Size() int64 {
	return rw.size
}
