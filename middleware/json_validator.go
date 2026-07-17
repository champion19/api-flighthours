package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/champion19/api-flighthours/platform/logger"
	json_schema "github.com/champion19/api-flighthours/platform/schema"
	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
)

type Builder struct {
	Validators *json_schema.Validators
	isLogin    bool
}

func NewMiddlewareValidator(validators *json_schema.Validators) *Builder {
	return &Builder{
		Validators: validators,
	}
}

// fieldNameMapping maps JSON field names to English labels for user-friendly error messages
var fieldNameMapping = map[string]string{
	// Person / Auth
	"name":                 "Name",
	"email":                "Email",
	"password":             "Password",
	"identificationnumber": "Identification Number",
	"identificationNumber": "Identification Number",
	"role":                 "Role",
	"token":                "Token",
	"current_password":     "Current Password",
	"new_password":         "New Password",
	"confirm_password":     "Confirm Password",

	// Employee / Airline Employee
	"active":     "Active",
	"start_date": "Start Date",
	"end_date":   "End Date",
	"airline_id": "Airline",
	"bp":         "BP Number",

	// Tail Number
	"tail_number":       "Tail Number",
	"aircraft_model_id": "Aircraft Model",

	// Daily Logbook
	"log_date":  "Log Date",
	"book_page": "Book Page",

	// Daily Logbook Detail (Flight)
	"airline_route_id":  "Airline Route",
	"tail_number_id":    "Aircraft Tail Number",
	"flight_number":     "Flight Number",
	"flight_real_date":  "Flight Date",
	"flight_type":       "Flight Type",
	"crew_role":         "Crew Role",
	"pilot_role":        "Pilot Role",
	"approach_category": "Approach Category",
	"approach_subtype":  "Approach Subtype",
	"autoland":          "Autoland",
	"companion_name":    "Companion Name",
	"passengers":        "Passengers",
	"out_time":          "Out Time",
	"takeoff_time":      "Takeoff Time",
	"landing_time":      "Landing Time",
	"in_time":           "In Time",
	"block_time":        "Block Time",
	"air_time":          "Air Time",

	// Message Management
	"code":     "Code",
	"type":     "Type",
	"title":    "Title",
	"content":  "Content",
	"module":   "Module",
	"category": "Category",
}

// translateFieldNames converts technical field names to English labels
func translateFieldNames(fields []string) []string {
	translated := make([]string, len(fields))
	for i, field := range fields {
		if label, exists := fieldNameMapping[field]; exists {
			translated[i] = label
		} else {
			translated[i] = field // Keep original if no mapping
		}
	}
	return translated
}

func (b *Builder) WithValidateRegister() gin.HandlerFunc {
	b.isLogin = false
	return b.jsonValidator(b.Validators.RegisterValidator)
}
func (b *Builder) WithValidateLogin() gin.HandlerFunc {
	b.isLogin = true
	return b.jsonValidator(b.Validators.LoginValidator)
}
func (b *Builder) WithValidateVerifyEmail() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.VerifyEmailValidator)
}
func (b *Builder) WithValidateMessage() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.MessageValidator)
}
func (b *Builder) WithValidateResendVerificationEmail() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ResendVerificationEmailValidator)
}
func (b *Builder) WithValidatePasswordResetRequest() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.PasswordResetRequestValidator)
}
func (b *Builder) WithValidateUpdatePassword() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdatePasswordValidator)
}
func (b *Builder) WithValidateUpdateEmployee() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateEmployeeValidator)
}
func (b *Builder) WithValidateChangePassword() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.ChangePasswordValidator)
}
func (b *Builder) WithValidateAddAirlineEmployee() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.AddAirlineEmployeeValidator)
}
func (b *Builder) WithValidateUpdateAirlineEmployee() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateAirlineEmployeeValidator)
}
func (b *Builder) WithValidateCreateTailNumber() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.CreateTailNumberValidator)
}
func (b *Builder) WithValidateUpdateTailNumber() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateTailNumberValidator)
}
func (b *Builder) WithValidateCreateDailyLogbookDetail() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.CreateDailyLogbookDetailValidator)
}
func (b *Builder) WithValidateUpdateDailyLogbookDetail() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateDailyLogbookDetailValidator)
}
func (b *Builder) WithValidateCreateDailyLogbook() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.CreateDailyLogbookValidator)
}
func (b *Builder) WithValidateUpdateDailyLogbook() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.UpdateDailyLogbookValidator)
}
func (b *Builder) WithValidateRefreshToken() gin.HandlerFunc {
	return b.jsonValidator(b.Validators.RefreshTokenValidator)
}

func (b *Builder) jsonValidator(schema *jsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := GetRequestID(c)
		log := log.WithTraceID(traceId)

		data, err := readAndParseBody(c, log)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		result := schema.Validate(data)
		if !result.IsValid() {
			handleValidationFailure(c, log, result)
			return
		}

		log.Debug(logger.LogMiddlewareValidationOK, "path", c.Request.URL.Path)
		c.Next()
	}
}

func readAndParseBody(c *gin.Context, log logger.Logger) (map[string]interface{}, error) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(logger.LogMiddlewareBodyReadError, "error", err, "path", c.Request.URL.Path)
		return nil, json_schema.ErrBodyReadFailed
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Error(logger.LogMiddlewareJSONParseError, "error", err, "path", c.Request.URL.Path)
		return nil, json_schema.ErrBadRequest
	}

	return data, nil
}

func handleValidationFailure(c *gin.Context, log logger.Logger, result *jsonschema.EvaluationResult) {
	fieldNames := extractFieldNames(result.Errors)
	validationError := classifyValidationError(fieldNames, result.Errors)

	if len(fieldNames) > 0 {
		c.Set("validation_fields", translateFieldNames(fieldNames))
	}

	log.Warn(logger.LogMiddlewareValidationFailed, "path", c.Request.URL.Path, "fields", fieldNames)
	c.Error(validationError)
	c.Abort()
}

// ============================================
// jsonValidator helpers (extracted to reduce cognitive complexity)
// ============================================

// extractFieldNames parses validation error params to collect field names.
// Handles both "properties" (plural, multiple fields) and "property" (singular).
func extractFieldNames(errors map[string]*jsonschema.EvaluationError) []string {
	var fieldNames []string
	for _, validationError := range errors {
		if validationError.Params == nil {
			continue
		}
		// Try "properties" (plural) first - for multiple fields
		if properties, exists := validationError.Params["properties"]; exists {
			fieldNames = append(fieldNames, parsePropertiesParam(properties)...)
		} else if property, exists := validationError.Params["property"]; exists {
			fieldNames = append(fieldNames, parsePropertiesParam(property)...)
		}
	}
	return fieldNames
}

// parsePropertiesParam splits a comma-separated property value into trimmed field names.
func parsePropertiesParam(value interface{}) []string {
	str := fmt.Sprintf("%v", value)
	var result []string
	for _, field := range strings.Split(str, ",") {
		field = strings.TrimSpace(field)
		field = strings.Trim(field, "'\"")
		if field != "" {
			result = append(result, field)
		}
	}
	return result
}

// classifyValidationError determines the specific validation error type
// based on field count and the first error code.
func classifyValidationError(fieldNames []string, errors map[string]*jsonschema.EvaluationError) error {
	if len(fieldNames) > 1 {
		return json_schema.ErrMultipleFields
	}

	// Single field error - determine specific error type
	var firstError *jsonschema.EvaluationError
	for _, err := range errors {
		firstError = err
		break
	}

	if firstError == nil {
		return json_schema.ErrValidationFailed
	}

	switch firstError.Code {
	case "property_mismatch":
		return json_schema.ErrFieldPropertyMismatch
	case "required":
		return json_schema.ErrFieldRequired
	case "type":
		return json_schema.ErrFieldTypeInvalid
	default:
		return json_schema.ErrValidationFailed
	}
}
