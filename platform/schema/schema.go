package schema

import (
	"io"
	"os"
	"path/filepath"

	"github.com/champion19/api-flighthours/tools/utils"
	"github.com/kaptinlin/jsonschema"
)

type Validators struct {
	FileReader                       FileReaderInterface
	RegisterValidator                *jsonschema.Schema
	MessageValidator                 *jsonschema.Schema
	ResendVerificationEmailValidator *jsonschema.Schema
	PasswordResetRequestValidator    *jsonschema.Schema
	UpdatePasswordValidator          *jsonschema.Schema
	UpdateEmployeeValidator          *jsonschema.Schema
	ChangePasswordValidator          *jsonschema.Schema
}

type FileReaderInterface interface {
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

func NewValidator(fileReader FileReaderInterface) (*Validators, error) {
	validator := &Validators{
		FileReader: fileReader,
	}

	register, err := validator.createSchema("register_person_schema.json")
	if err != nil {
		return nil, err
	}
	message, err := validator.createSchema("message_schema.json")
	if err != nil {
		return nil, err
	}
	resendVerificationEmail, err := validator.createSchema("resend_verification_email_schema.json")
	if err != nil {
		return nil, err
	}
	passwordResetRequest, err := validator.createSchema("password_reset_request_schema.json")
	if err != nil {
		return nil, err
	}
	updatePassword, err := validator.createSchema("update_password_schema.json")
	if err != nil {
		return nil, err
	}
	updateEmployee, err := validator.createSchema("update_employee_schema.json")
	if err != nil {
		return nil, err
	}
	changePassword, err := validator.createSchema("change_password_schema.json")
	if err != nil {
		return nil, err
	}

	validator.RegisterValidator = register
	validator.MessageValidator = message
	validator.ResendVerificationEmailValidator = resendVerificationEmail
	validator.PasswordResetRequestValidator = passwordResetRequest
	validator.UpdatePasswordValidator = updatePassword
	validator.UpdateEmployeeValidator = updateEmployee
	validator.ChangePasswordValidator = changePassword

	return validator, nil

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
