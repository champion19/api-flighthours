package prometheus

// ============================================
// PROMETHEUS METRICS CONSTANTS
// ============================================
// Centralized constants for Prometheus metrics
// to avoid hardcoded strings across the codebase

// ============================================
// METRIC NAMES
// ============================================
const (
	// HTTP Metrics
	MetricHTTPRequestsTotal   = "flighthours_http_requests_total"
	MetricHTTPRequestDuration = "flighthours_http_request_duration_seconds"
	MetricHTTPErrorsTotal     = "flighthours_http_errors_total"

	// Business Metrics
	MetricUserRegistrationsTotal = "flighthours_user_registrations_total"
	MetricMessagesCreatedTotal   = "flighthours_messages_created_total"
)

// ============================================
// METRIC DESCRIPTIONS
// ============================================
const (
	// HTTP Metrics Descriptions
	MetricHTTPRequestsTotalHelp   = "Total number of HTTP requests processed by Motogo Backend"
	MetricHTTPRequestDurationHelp = "Duration of HTTP requests in seconds"
	MetricHTTPErrorsTotalHelp     = "Total number of HTTP errors (status >= 400)"

	// Business Metrics Descriptions
	MetricUserRegistrationsTotalHelp = "Total number of user registrations"
	MetricMessagesCreatedTotalHelp   = "Total number of messages created"
)

// ============================================
// LABEL NAMES
// ============================================
const (
	LabelMethod    = "method"
	LabelEndpoint  = "endpoint"
	LabelStatus    = "status"
	LabelErrorType = "error_type"
	LabelModule    = "module"
	LabelType      = "type"
)

// ============================================
// ERROR TYPES
// ============================================
const (
	ErrorTypeServerError = "server_error"
	ErrorTypeNotFound    = "not_found"
	ErrorTypeAuthError   = "auth_error"
	ErrorTypeBadRequest  = "bad_request"
	ErrorTypeValidation  = "validation_error"
	ErrorTypeConflict    = "conflict"
	ErrorTypeRateLimit   = "rate_limit_error"
	ErrorTypeClientError = "client_error"
)
