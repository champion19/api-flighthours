package schema

import (
	"io"
	"os"
	"path/filepath"

	"github.com/champion19/api-flighthours/tools/utils"
	"github.com/kaptinlin/jsonschema"
)

type Validators struct {
	FileReader                        SchemaReader
	RegisterValidator                 *jsonschema.Schema
	LoginValidator                    *jsonschema.Schema
	MessageValidator                  *jsonschema.Schema
	ResendVerificationEmailValidator  *jsonschema.Schema
	PasswordResetRequestValidator     *jsonschema.Schema
	UpdatePasswordValidator           *jsonschema.Schema
	UpdateEmployeeValidator           *jsonschema.Schema
	ChangePasswordValidator           *jsonschema.Schema
	VerifyEmailValidator              *jsonschema.Schema
	AddAirlineEmployeeValidator       *jsonschema.Schema
	UpdateAirlineEmployeeValidator    *jsonschema.Schema
	CreateTailNumberValidator         *jsonschema.Schema
	UpdateTailNumberValidator         *jsonschema.Schema
	CreateDailyLogbookDetailValidator *jsonschema.Schema
	UpdateDailyLogbookDetailValidator *jsonschema.Schema
	CreateDailyLogbookValidator       *jsonschema.Schema
	UpdateDailyLogbookValidator       *jsonschema.Schema
	RefreshTokenValidator             *jsonschema.Schema
}

// SchemaReader reads JSON schema files from the filesystem.
type SchemaReader interface {
	ReadJsonSchema(resourcePath string) ([]byte, error)
}

type DefaultFileReader struct{}

func (f *DefaultFileReader) ReadJsonSchema(resourcePath string) ([]byte, error) {

	root, err := utils.FindModuleRoot()

	if err != nil {
		return nil, err
	}

	data, err := os.Open(filepath.Join(root, "platform/schema/json_schema", resourcePath))
	if err != nil {
		return nil, err
	}
	defer data.Close()

	return io.ReadAll(data)

}

func NewValidator(fileReader SchemaReader) (*Validators, error) {
	v := &Validators{FileReader: fileReader}

	type schemaEntry struct {
		file   string
		target **jsonschema.Schema
	}

	entries := []schemaEntry{
		{"register_person_schema.json", &v.RegisterValidator},
		{"login_schema.json", &v.LoginValidator},
		{"message_schema.json", &v.MessageValidator},
		{"resend_verification_email_schema.json", &v.ResendVerificationEmailValidator},
		{"password_reset_request_schema.json", &v.PasswordResetRequestValidator},
		{"update_password_schema.json", &v.UpdatePasswordValidator},
		{"update_employee_schema.json", &v.UpdateEmployeeValidator},
		{"change_password_schema.json", &v.ChangePasswordValidator},
		{"verify_email_schema.json", &v.VerifyEmailValidator},
		{"add_airline_employee_schema.json", &v.AddAirlineEmployeeValidator},
		{"update_airline_employee_schema.json", &v.UpdateAirlineEmployeeValidator},
		{"create_tail_number_schema.json", &v.CreateTailNumberValidator},
		{"update_tail_number_schema.json", &v.UpdateTailNumberValidator},
		{"create_daily_logbook_detail_schema.json", &v.CreateDailyLogbookDetailValidator},
		{"update_daily_logbook_detail_schema.json", &v.UpdateDailyLogbookDetailValidator},
		{"create_daily_logbook_schema.json", &v.CreateDailyLogbookValidator},
		{"update_daily_logbook_schema.json", &v.UpdateDailyLogbookValidator},
		{"refresh_token_schema.json", &v.RefreshTokenValidator},
	}

	for _, e := range entries {
		s, err := v.createSchema(e.file)
		if err != nil {
			return nil, err
		}
		*e.target = s
	}

	return v, nil
}

func (v *Validators) createSchema(resourcePath string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat = true
	schemaJSON, err := v.FileReader.ReadJsonSchema(resourcePath)
	if err != nil {
		return nil, ErrSchemaReadFailed
	}

	if schemaJSON == nil {
		return nil, ErrSchemaEmpty
	}

	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		return nil, ErrSchemaCompileFailed
	}

	return schema, nil
}

func (v *Validators) ValidateRegister(data interface{}) error {
	if v.RegisterValidator == nil {
		return ErrSchemaEmpty
	}

	result := v.RegisterValidator.Validate(data)
	if !result.IsValid() {

		var errorMessages []string
		for _, err := range result.Errors {
			errorMessages = append(errorMessages, err.Message)
		}

		if len(errorMessages) > 0 {
			return ErrValidationFailed
		}
	}

	return nil
}

func (v *Validators) ValidateUpdateEmployee(data interface{}) error {
	if v.UpdateEmployeeValidator == nil {
		return ErrSchemaEmpty
	}

	result := v.UpdateEmployeeValidator.Validate(data)
	if !result.IsValid() {

		var errorMessages []string
		for _, err := range result.Errors {
			errorMessages = append(errorMessages, err.Message)
		}

		if len(errorMessages) > 0 {
			return ErrValidationFailed
		}
	}

	return nil
}
