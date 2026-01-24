package middleware

import (
	"strconv"
	"time"

	promConstants "github.com/champion19/flighthours-api/platform/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTP Request counter
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: promConstants.MetricHTTPRequestsTotal,
			Help: promConstants.MetricHTTPRequestsTotalHelp,
		},
		[]string{promConstants.LabelMethod, promConstants.LabelEndpoint, promConstants.LabelStatus},
	)

	// HTTP Request duration histogram
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    promConstants.MetricHTTPRequestDuration,
			Help:    promConstants.MetricHTTPRequestDurationHelp,
			Buckets: prometheus.DefBuckets, // [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
		},
		[]string{promConstants.LabelMethod, promConstants.LabelEndpoint, promConstants.LabelStatus},
	)

	// Error counter
	httpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: promConstants.MetricHTTPErrorsTotal,
			Help: promConstants.MetricHTTPErrorsTotalHelp,
		},
		[]string{promConstants.LabelMethod, promConstants.LabelEndpoint, promConstants.LabelStatus, promConstants.LabelErrorType},
	)

	// Business metrics - User registrations
	userRegistrationsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: promConstants.MetricUserRegistrationsTotal,
			Help: promConstants.MetricUserRegistrationsTotalHelp,
		},
	)

	// Business metrics - Messages created
	messagesCreatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: promConstants.MetricMessagesCreatedTotal,
			Help: promConstants.MetricMessagesCreatedTotalHelp,
		},
		[]string{promConstants.LabelModule, promConstants.LabelType},
	)
)

func PrometheusInit() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpErrorsTotal)
	prometheus.MustRegister(userRegistrationsTotal)
	prometheus.MustRegister(messagesCreatedTotal)
}

// TrackMetrics is a Gin middleware that tracks HTTP metrics
func TrackMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Process request
		c.Next()

		// Record metrics after request processing
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Record request count
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()

		// Record request duration
		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)

		// Record errors if status >= 400
		if c.Writer.Status() >= 400 {
			errorType := getErrorType(c.Writer.Status())
			httpErrorsTotal.WithLabelValues(method, path, status, errorType).Inc()
		}
	}
}

// getErrorType classifies HTTP status codes into error types
func getErrorType(status int) string {
	switch {
	case status >= 500:
		return promConstants.ErrorTypeServerError
	case status == 404:
		return promConstants.ErrorTypeNotFound
	case status == 401 || status == 403:
		return promConstants.ErrorTypeAuthError
	case status == 400:
		return promConstants.ErrorTypeBadRequest
	case status == 422:
		return promConstants.ErrorTypeValidation
	case status == 409:
		return promConstants.ErrorTypeConflict
	case status == 429:
		return promConstants.ErrorTypeRateLimit
	default:
		return promConstants.ErrorTypeClientError
	}
}

// RecordEmployeeRegistration increments the employee registration counter
func RecordEmployeeRegistration() {
	userRegistrationsTotal.Inc()
}

// RecordMessageCreated increments the message creation counter
func RecordMessageCreated(module, msgType string) {
	messagesCreatedTotal.WithLabelValues(module, msgType).Inc()
}

// PrometheusMiddleware captura métricas de cada request HTTP
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Process request
		c.Next()

		// Record metrics after request processing
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// Record request count
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()

		// Record request duration
		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)

		// Record errors if status >= 400
		if c.Writer.Status() >= 400 {
			errorType := getErrorType(c.Writer.Status())
			httpErrorsTotal.WithLabelValues(method, path, status, errorType).Inc()
		}
	}
}
