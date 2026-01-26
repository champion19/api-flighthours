package middleware

import (
	"strconv"
	"time"

	promConstants "github.com/champion19/flighthours-api/platform/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: promConstants.MetricHTTPRequestsTotal,
			Help: promConstants.MetricHTTPRequestsTotalHelp,
		},
		[]string{promConstants.LabelMethod, promConstants.LabelEndpoint, promConstants.LabelStatus},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    promConstants.MetricHTTPRequestDuration,
			Help:    promConstants.MetricHTTPRequestDurationHelp,
			Buckets: prometheus.DefBuckets,
		},
		[]string{promConstants.LabelMethod, promConstants.LabelEndpoint, promConstants.LabelStatus},
	)

	httpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: promConstants.MetricHTTPErrorsTotal,
			Help: promConstants.MetricHTTPErrorsTotalHelp,
		},
		[]string{promConstants.LabelMethod, promConstants.LabelEndpoint, promConstants.LabelStatus, promConstants.LabelErrorType},
	)

	userRegistrationsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: promConstants.MetricUserRegistrationsTotal,
			Help: promConstants.MetricUserRegistrationsTotalHelp,
		},
	)

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

func TrackMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()

		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)

		if c.Writer.Status() >= 400 {
			errorType := getErrorType(c.Writer.Status())
			httpErrorsTotal.WithLabelValues(method, path, status, errorType).Inc()
		}
	}
}

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

func RecordEmployeeRegistration() {
	userRegistrationsTotal.Inc()
}
func RecordMessageCreated(module, msgType string) {
	messagesCreatedTotal.WithLabelValues(module, msgType).Inc()
}


func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()

		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)

		if c.Writer.Status() >= 400 {
			errorType := getErrorType(c.Writer.Status())
			httpErrorsTotal.WithLabelValues(method, path, status, errorType).Inc()
		}
	}
}
