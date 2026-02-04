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
