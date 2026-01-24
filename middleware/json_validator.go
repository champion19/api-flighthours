package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/champion19/flighthours-api/platform/logger"
	json_schema "github.com/champion19/flighthours-api/platform/schema"
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

func (b *Builder) WithValidateRegister() gin.HandlerFunc {
	b.isLogin = false
	return b.jsonValidator(b.Validators.RegisterValidator)
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


func (b *Builder) jsonValidator(schema *jsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {

		traceId := GetRequestID(c)
		log := log.WithTraceID(traceId)

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if log != nil {
				log.Error(logger.LogMiddlewareBodyReadError, "error", err, "path", c.Request.URL.Path)
			}
			c.Error(json_schema.ErrBodyReadFailed)
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		var data map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			if log != nil {
				log.Error(logger.LogMiddlewareJSONParseError, "error", err, "path", c.Request.URL.Path)
			}
			c.Error(json_schema.ErrBadRequest)
			c.Abort()
			return
		}
		result := schema.Validate(data)
		if !result.IsValid() {
			var fieldNames []string

			for _, validationError := range result.Errors {
				if validationError.Params != nil {
					if properties, exists := validationError.Params["properties"]; exists {
						propertiesStr := fmt.Sprintf("%v", properties)
						fields := strings.Split(propertiesStr, ",")
						for _, field := range fields {
							field = strings.TrimSpace(field)
							field = strings.Trim(field, "'\"")
							if field != "" {
								fieldNames = append(fieldNames, field)
							}
						}
					} else if property, exists := validationError.Params["property"]; exists {
						propertyName := fmt.Sprintf("%v", property)
						propertyName = strings.Trim(propertyName, "'\"")
						if propertyName != "" {
							fieldNames = append(fieldNames, propertyName)
						}
					}
				}
			}

			var validationError error

			// If multiple fields failed, use specific multiple fields error
			if len(fieldNames) > 1 {
				validationError = json_schema.ErrMultipleFields
			} else {
				// Single field error - determine specific error type
				var firstError *jsonschema.EvaluationError
				for _, err := range result.Errors {
					firstError = err
					break
				}

				if firstError != nil {
					switch firstError.Code {
					case "property_mismatch":
						validationError = json_schema.ErrFieldPropertyMismatch
					case "required":
						validationError = json_schema.ErrFieldRequired
					case "type":
						validationError = json_schema.ErrFieldTypeInvalid
					default:
						validationError = json_schema.ErrValidationFailed
					}
				} else {
					validationError = json_schema.ErrValidationFailed
				}
			}

			// Store field names in context for error_handler to use in message parameters
			if len(fieldNames) > 0 {
				c.Set("validation_fields", fieldNames)
			}

			if log != nil {
				log.Warn(logger.LogMiddlewareValidationFailed, "path", c.Request.URL.Path, "fields", fieldNames)
			}
			c.Error(validationError)
			c.Abort()
			return
		}

		if log != nil {
			log.Debug(logger.LogMiddlewareValidationOK, "path", c.Request.URL.Path)
		}

		c.Next()
	}
}
