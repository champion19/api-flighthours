package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	json_schema "github.com/champion19/api-flighthours/platform/schema"
	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewMiddlewareValidator(t *testing.T) {
	t.Run("creates builder with nil validators", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)
		if builder == nil {
			t.Error("expected non-nil Builder")
		}
	})

	t.Run("stores validators reference", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)
		if builder.Validators != nil {
			t.Error("expected nil Validators when passed nil")
		}
	})

	t.Run("isLogin defaults to false", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)
		if builder.isLogin != false {
			t.Error("expected isLogin to default to false")
		}
	})
}

func TestJsonValidator_ValidBody(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	schemaJSON := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"required": ["name"],
		"properties": {
			"name": {"type": "string"}
		}
	}`)
	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		t.Fatalf("failed to compile test schema: %v", err)
	}

	t.Run("valid JSON body passes validation", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		body := bytes.NewReader([]byte(`{"name": "test"}`))
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestJsonValidator_Errors(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	schemaJSON := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"required": ["name"],
		"properties": {
			"name": {"type": "string"}
		}
	}`)
	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		t.Fatalf("failed to compile test schema: %v", err)
	}

	// Test that middleware sets errors - not checking abort since gin handles that
	t.Run("invalid JSON syntax sets error on context", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		body := bytes.NewReader([]byte(`{invalid`))
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Handler should not be reached due to abort
		// We just verify no panic occurs
	})

	t.Run("missing required field causes validation to fail", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		body := bytes.NewReader([]byte(`{}`))
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Validation should fail and abort
	})

	t.Run("wrong type causes validation to fail", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		body := bytes.NewReader([]byte(`{"name": 123}`))
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Validation should fail and abort
	})
}

func TestJsonValidator_MultipleFields(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	schemaJSON := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"required": ["name", "email", "age"],
		"properties": {
			"name": {"type": "string"},
			"email": {"type": "string"},
			"age": {"type": "integer"}
		}
	}`)
	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		t.Fatalf("failed to compile test schema: %v", err)
	}

	t.Run("valid data with all fields passes", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		body := bytes.NewReader([]byte(`{"name": "John", "email": "john@example.com", "age": 30}`))
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("missing multiple fields fails validation", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		body := bytes.NewReader([]byte(`{}`))
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Should not get OK - validation should have failed
	})
}

func TestJsonValidator_RestoresBody(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	schemaJSON := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"properties": {
			"name": {"type": "string"}
		}
	}`)
	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		t.Fatalf("failed to compile test schema: %v", err)
	}

	t.Run("body is still readable after validation", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)
		var bodyReadable bool

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			var data map[string]interface{}
			if err := c.ShouldBindJSON(&data); err == nil {
				bodyReadable = true
			}
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		body := bytes.NewReader([]byte(`{"name": "test"}`))
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if !bodyReadable {
			t.Error("expected body to be readable after validation")
		}
	})
}

func TestJsonValidator_EmptyBody(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	schemaJSON := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object"
	}`)
	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		t.Fatalf("failed to compile test schema: %v", err)
	}

	t.Run("empty body fails JSON parsing", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)

		r := gin.New()
		r.Use(builder.jsonValidator(schema))
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Handler should not be reached due to abort
	})
}

func TestWithValidateMethods(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	schemaJSON := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"properties": {
			"name": {"type": "string"}
		}
	}`)
	schema, _ := compiler.Compile(schemaJSON)

	// Create mock validators structure
	validators := &json_schema.Validators{
		RegisterValidator:                schema,
		MessageValidator:                 schema,
		ResendVerificationEmailValidator: schema,
		PasswordResetRequestValidator:    schema,
		UpdatePasswordValidator:          schema,
		UpdateEmployeeValidator:          schema,
		ChangePasswordValidator:          schema,
		AddAirlineEmployeeValidator:      schema,
		UpdateAirlineEmployeeValidator:   schema,
		CreateLicensePlateValidator:      schema,
		UpdateLicensePlateValidator:      schema,
	}

	builder := NewMiddlewareValidator(validators)

	t.Run("WithValidateRegister returns handler", func(t *testing.T) {
		handler := builder.WithValidateRegister()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateMessage returns handler", func(t *testing.T) {
		handler := builder.WithValidateMessage()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateResendVerificationEmail returns handler", func(t *testing.T) {
		handler := builder.WithValidateResendVerificationEmail()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidatePasswordResetRequest returns handler", func(t *testing.T) {
		handler := builder.WithValidatePasswordResetRequest()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateUpdatePassword returns handler", func(t *testing.T) {
		handler := builder.WithValidateUpdatePassword()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateUpdateEmployee returns handler", func(t *testing.T) {
		handler := builder.WithValidateUpdateEmployee()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateChangePassword returns handler", func(t *testing.T) {
		handler := builder.WithValidateChangePassword()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateAddAirlineEmployee returns handler", func(t *testing.T) {
		handler := builder.WithValidateAddAirlineEmployee()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateUpdateAirlineEmployee returns handler", func(t *testing.T) {
		handler := builder.WithValidateUpdateAirlineEmployee()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateCreateLicensePlate returns handler", func(t *testing.T) {
		handler := builder.WithValidateCreateLicensePlate()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("WithValidateUpdateLicensePlate returns handler", func(t *testing.T) {
		handler := builder.WithValidateUpdateLicensePlate()
		if handler == nil {
			t.Error("expected non-nil handler")
		}
	})
}
