package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSONFormat           = errors.New("invalid JSON format")
	ErrUnmarshalBody               = errors.New("failed to process request body")
	ErrSchemaValidation            = errors.New("validation failed")
	ErrInternalServer              = errors.New("internal server error")
	ErrModuleRootNotFound          = errors.New("could not find module root")
	ErrSchemaFileNotFound          = errors.New("schema file not found")
	ErrSchemaFileRead              = errors.New("failed to read schema file")
	ErrSchemaCompilation           = errors.New("failed to compile JSON schema")
	ErrSchemaEmpty                 = errors.New("JSON schema is empty or null")
	ErrValidatorInitFailed         = errors.New("validator initialization failed")
	ErrValidationUserFailed        = errors.New("user validation failed")
	ErrValidationUserNotFound      = errors.New("user not found")
	ErrValidationUserAlreadyExists = errors.New("user already exists")
	ErrDBQueryFailed               = errors.New("database query failed")
)

func ValidateError(c *gin.Context, err error, details interface{}, statusCode int) {
	if details == nil {
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	detailsMap, ok := details.(map[string]interface{})
	if !ok {
		c.JSON(statusCode, gin.H{
			"error":   err.Error(),
			"details": details,
		})
		c.Abort()
		return
	}

	fieldErrors := extractFieldErrors(detailsMap)

	c.JSON(statusCode, gin.H{
		"error":   err.Error(),
		"invalid": fieldErrors,
	})
	c.Abort()
}

func extractFieldErrors(detailsMap map[string]interface{}) map[string]string {
	fieldErrors := make(map[string]string)

	detailsList, exists := detailsMap["details"].([]interface{})
	if !exists {
		return fieldErrors
	}

	for _, item := range detailsList {
		fieldDetail, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		parseFieldDetail(fieldDetail, fieldErrors)
	}

	return fieldErrors
}

func parseFieldDetail(fieldDetail map[string]interface{}, fieldErrors map[string]string) {
	valid, exists := fieldDetail["valid"].(bool)
	if !exists || valid {
		return
	}

	path := extractInstancePath(fieldDetail)
	if path == "" {
		return
	}

	errorsMap, exists := fieldDetail["errors"].(map[string]interface{})
	if !exists {
		return
	}

	for _, msg := range errorsMap {
		if strMsg, ok := msg.(string); ok {
			fieldErrors[path] = strMsg
		}
	}
}

func extractInstancePath(fieldDetail map[string]interface{}) string {
	instanceLoc, exists := fieldDetail["instanceLocation"].(string)
	if !exists {
		return ""
	}
	if len(instanceLoc) > 0 && instanceLoc[0] == '/' {
		return instanceLoc[1:]
	}
	return instanceLoc
}
