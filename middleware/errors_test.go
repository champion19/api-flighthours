package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestValidateError_WithNilDetails(t *testing.T) {
	t.Run("returns error without details", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ValidateError(c, ErrSchemaValidation, nil, http.StatusBadRequest)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
		if c.IsAborted() != true {
			t.Error("expected context to be aborted")
		}
	})
}

func TestValidateError_WithNonMapDetails(t *testing.T) {
	t.Run("returns error with raw details", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		details := "some string details"
		ValidateError(c, ErrSchemaValidation, details, http.StatusUnprocessableEntity)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
		}
	})
}

func TestValidateError_WithMapDetails(t *testing.T) {
	t.Run("parses field errors from details map", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		details := map[string]interface{}{
			"details": []interface{}{
				map[string]interface{}{
					"valid":            false,
					"instanceLocation": "/email",
					"errors": map[string]interface{}{
						"format": "must be a valid email",
					},
				},
			},
		}

		ValidateError(c, ErrSchemaValidation, details, http.StatusBadRequest)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("handles details with valid true", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		details := map[string]interface{}{
			"details": []interface{}{
				map[string]interface{}{
					"valid":            true,
					"instanceLocation": "/name",
				},
			},
		}

		ValidateError(c, ErrSchemaValidation, details, http.StatusBadRequest)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("handles details with invalid item format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		details := map[string]interface{}{
			"details": []interface{}{
				"not a map",
			},
		}

		ValidateError(c, ErrSchemaValidation, details, http.StatusBadRequest)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestMiddlewareErrors(t *testing.T) {
	t.Run("error messages are not empty", func(t *testing.T) {
		errors := []error{
			ErrInvalidJSONFormat,
			ErrUnmarshalBody,
			ErrSchemaValidation,
			ErrInternalServer,
			ErrModuleRootNotFound,
			ErrSchemaFileNotFound,
			ErrSchemaFileRead,
			ErrSchemaCompilation,
			ErrSchemaEmpty,
			ErrValidatorInitFailed,
			ErrValidationUserFailed,
			ErrValidationUserNotFound,
			ErrValidationUserAlreadyExists,
			ErrDBQueryFailed,
		}

		for _, err := range errors {
			if err.Error() == "" {
				t.Errorf("error should have non-empty message: %v", err)
			}
		}
	})
}
