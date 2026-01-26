package prometheus




const (
	MetricHTTPRequestsTotal   = "flighthours_http_requests_total"
	MetricHTTPRequestDuration = "flighthours_http_request_duration_seconds"
	MetricHTTPErrorsTotal     = "flighthours_http_errors_total"
	MetricUserRegistrationsTotal = "flighthours_user_registrations_total"
	MetricMessagesCreatedTotal   = "flighthours_messages_created_total"
)


const (
	MetricHTTPRequestsTotalHelp   = "Total number of HTTP requests processed by Motogo Backend"
	MetricHTTPRequestDurationHelp = "Duration of HTTP requests in seconds"
	MetricHTTPErrorsTotalHelp     = "Total number of HTTP errors (status >= 400)"


	MetricUserRegistrationsTotalHelp = "Total number of user registrations"
	MetricMessagesCreatedTotalHelp   = "Total number of messages created"
)


const (
	LabelMethod    = "method"
	LabelEndpoint  = "endpoint"
	LabelStatus    = "status"
	LabelErrorType = "error_type"
	LabelModule    = "module"
	LabelType      = "type"
)


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
