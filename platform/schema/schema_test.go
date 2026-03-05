package schema

import (
	"errors"
	"testing"
)

// Mock file reader for testing
type mockFileReader struct {
	data    []byte
	readErr error
}

func (m *mockFileReader) ReadJsonSchema(resourcePath string) ([]byte, error) {
	if m.readErr != nil {
		return nil, m.readErr
	}
	return m.data, nil
}

func TestValidators_createSchema(t *testing.T) {
	t.Run("returns error when file read fails", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{readErr: errors.New("file not found")},
		}

		_, err := v.createSchema("test.json")

		if err == nil {
			t.Error("expected error when file read fails")
		}
		if !errors.Is(err, ErrSchemaReadFailed) {
			t.Errorf("expected ErrSchemaReadFailed, got %v", err)
		}
	})

	t.Run("returns error when schema is empty", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: nil},
		}

		_, err := v.createSchema("test.json")

		if err == nil {
			t.Error("expected error when schema is empty")
		}
		if !errors.Is(err, ErrSchemaEmpty) {
			t.Errorf("expected ErrSchemaEmpty, got %v", err)
		}
	})

	t.Run("returns error when schema is invalid JSON", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: []byte("not valid json")},
		}

		_, err := v.createSchema("test.json")

		if err == nil {
			t.Error("expected error when schema is invalid")
		}
		if !errors.Is(err, ErrSchemaCompileFailed) {
			t.Errorf("expected ErrSchemaCompileFailed, got %v", err)
		}
	})

	t.Run("creates schema successfully", func(t *testing.T) {
		validSchema := []byte(`{
			"$schema": "https://json-schema.org/draft/2020-12/schema",
			"type": "object",
			"properties": {
				"name": {"type": "string"}
			}
		}`)
		v := &Validators{
			FileReader: &mockFileReader{data: validSchema},
		}

		schema, err := v.createSchema("test.json")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if schema == nil {
			t.Error("expected non-nil schema")
		}
	})
}

func TestValidators_ValidateRegister(t *testing.T) {
	t.Run("returns error when validator is nil", func(t *testing.T) {
		v := &Validators{}

		err := v.ValidateRegister(map[string]interface{}{})

		if err == nil {
			t.Error("expected error when validator is nil")
		}
		if !errors.Is(err, ErrSchemaEmpty) {
			t.Errorf("expected ErrSchemaEmpty, got %v", err)
		}
	})
}

func TestValidators_ValidateUpdateEmployee(t *testing.T) {
	t.Run("returns error when validator is nil", func(t *testing.T) {
		v := &Validators{}

		err := v.ValidateUpdateEmployee(map[string]interface{}{})

		if err == nil {
			t.Error("expected error when validator is nil")
		}
		if !errors.Is(err, ErrSchemaEmpty) {
			t.Errorf("expected ErrSchemaEmpty, got %v", err)
		}
	})
}

func TestSchemaErrors(t *testing.T) {
	t.Run("error constants are defined", func(t *testing.T) {
		if ErrSchemaReadFailed == nil {
			t.Error("ErrSchemaReadFailed should not be nil")
		}
		if ErrSchemaEmpty == nil {
			t.Error("ErrSchemaEmpty should not be nil")
		}
		if ErrSchemaCompileFailed == nil {
			t.Error("ErrSchemaCompileFailed should not be nil")
		}
		if ErrValidationFailed == nil {
			t.Error("ErrValidationFailed should not be nil")
		}
	})
}

func TestValidators_ValidateRegister_WithData(t *testing.T) {
	// Valid schema for registration testing
	registerSchema := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"required": ["email", "password"],
		"properties": {
			"email": {"type": "string", "format": "email"},
			"password": {"type": "string", "minLength": 8}
		}
	}`)

	t.Run("returns nil for valid data", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: registerSchema},
		}
		schema, err := v.createSchema("test.json")
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
		v.RegisterValidator = schema

		err = v.ValidateRegister(map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		})

		if err != nil {
			t.Errorf("expected nil error for valid data, got %v", err)
		}
	})

	t.Run("returns error for invalid email format", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: registerSchema},
		}
		schema, err := v.createSchema("test.json")
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
		v.RegisterValidator = schema

		err = v.ValidateRegister(map[string]interface{}{
			"email":    "not-an-email",
			"password": "password123",
		})

		if err == nil {
			t.Error("expected error for invalid email format")
		}
	})

	t.Run("returns error for short password", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: registerSchema},
		}
		schema, err := v.createSchema("test.json")
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
		v.RegisterValidator = schema

		err = v.ValidateRegister(map[string]interface{}{
			"email":    "test@example.com",
			"password": "short",
		})

		if err == nil {
			t.Error("expected error for short password")
		}
	})

	t.Run("returns error for missing required fields", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: registerSchema},
		}
		schema, err := v.createSchema("test.json")
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
		v.RegisterValidator = schema

		err = v.ValidateRegister(map[string]interface{}{})

		if err == nil {
			t.Error("expected error for missing required fields")
		}
	})
}

func TestValidators_ValidateUpdateEmployee_WithData(t *testing.T) {
	// Valid schema for update employee testing
	updateSchema := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"properties": {
			"name": {"type": "string", "minLength": 1},
			"role": {"type": "string", "enum": ["admin", "pilot", "user"]}
		}
	}`)

	t.Run("returns nil for valid update data", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: updateSchema},
		}
		schema, err := v.createSchema("test.json")
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
		v.UpdateEmployeeValidator = schema

		err = v.ValidateUpdateEmployee(map[string]interface{}{
			"name": "John Doe",
			"role": "pilot",
		})

		if err != nil {
			t.Errorf("expected nil error for valid data, got %v", err)
		}
	})

	t.Run("returns error for invalid role (not in enum)", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: updateSchema},
		}
		schema, err := v.createSchema("test.json")
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
		v.UpdateEmployeeValidator = schema

		err = v.ValidateUpdateEmployee(map[string]interface{}{
			"name": "John Doe",
			"role": "invalid_role",
		})

		if err == nil {
			t.Error("expected error for invalid role")
		}
	})

	t.Run("returns nil for empty object (no required fields)", func(t *testing.T) {
		v := &Validators{
			FileReader: &mockFileReader{data: updateSchema},
		}
		schema, err := v.createSchema("test.json")
		if err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
		v.UpdateEmployeeValidator = schema

		err = v.ValidateUpdateEmployee(map[string]interface{}{})

		if err != nil {
			t.Errorf("expected nil error for empty object (no required fields), got %v", err)
		}
	})
}

func TestDefaultFileReader_Interface(t *testing.T) {
	t.Run("implements FileReaderInterface", func(t *testing.T) {
		var _ SchemaReader = (*DefaultFileReader)(nil)
	})
}

func TestNewValidator_Success(t *testing.T) {
	// Use a mock that returns valid schemas
	validSchema := []byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"properties": {
			"name": {"type": "string"}
		}
	}`)

	t.Run("creates validator with all schemas", func(t *testing.T) {
		mock := &mockFileReader{data: validSchema}
		v, err := NewValidator(mock)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v == nil {
			t.Fatal("expected non-nil validator")
		}
		if v.RegisterValidator == nil {
			t.Error("expected RegisterValidator to be set")
		}
		if v.MessageValidator == nil {
			t.Error("expected MessageValidator to be set")
		}
		if v.PasswordResetRequestValidator == nil {
			t.Error("expected PasswordResetRequestValidator to be set")
		}
	})

	t.Run("returns error when any schema fails", func(t *testing.T) {
		callCount := 0
		mockWithFailure := &mockConditionalFileReader{
			failOnCall: 3, // Fail on the third schema
			data:       validSchema,
		}
		_, err := NewValidator(mockWithFailure)
		_ = callCount // used in mock

		if err == nil {
			t.Error("expected error when schema creation fails")
		}
	})
}

// Mock that fails on a specific call
type mockConditionalFileReader struct {
	callCount  int
	failOnCall int
	data       []byte
}

func (m *mockConditionalFileReader) ReadJsonSchema(resourcePath string) ([]byte, error) {
	m.callCount++
	if m.callCount == m.failOnCall {
		return nil, errors.New("simulated failure")
	}
	return m.data, nil
}

func TestDefaultFileReader_ReadJsonSchema(t *testing.T) {
	t.Run("reads existing schema file", func(t *testing.T) {
		reader := &DefaultFileReader{}

		// This should work if the schema files exist
		data, err := reader.ReadJsonSchema("register_person_schema.json")

		if err != nil {
			t.Fatalf("unexpected error reading schema: %v", err)
		}
		if len(data) == 0 {
			t.Error("expected non-empty data")
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		reader := &DefaultFileReader{}

		_, err := reader.ReadJsonSchema("non_existent_schema.json")

		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})
}
