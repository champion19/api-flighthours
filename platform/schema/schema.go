package schema

import (
	"io"
	"os"
	"path/filepath"

	"github.com/champion19/api-flighthours/tools/utils"
	"github.com/kaptinlin/jsonschema"
)

type Validators struct {
	FileReader                        FileReaderInterface
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
	CreateLicensePlateValidator       *jsonschema.Schema
	UpdateLicensePlateValidator       *jsonschema.Schema
	CreateDailyLogbookDetailValidator *jsonschema.Schema
	UpdateDailyLogbookDetailValidator *jsonschema.Schema
	CreateDailyLogbookValidator       *jsonschema.Schema
	UpdateDailyLogbookValidator       *jsonschema.Schema
	RefreshTokenValidator             *jsonschema.Schema
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
	login, err := validator.createSchema("login_schema.json")
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
	verifyEmail, err := validator.createSchema("verify_email_schema.json")
	if err != nil {
		return nil, err
	}
	addAirlineEmployee, err := validator.createSchema("add_airline_employee_schema.json")
	if err != nil {
		return nil, err
	}
	updateAirlineEmployee, err := validator.createSchema("update_airline_employee_schema.json")
	if err != nil {
		return nil, err
	}
	createLicensePlate, err := validator.createSchema("create_license_plate_schema.json")
	if err != nil {
		return nil, err
	}
	updateLicensePlate, err := validator.createSchema("update_license_plate_schema.json")
	if err != nil {
		return nil, err
	}
	createDailyLogbookDetail, err := validator.createSchema("create_daily_logbook_detail_schema.json")
	if err != nil {
		return nil, err
	}
	updateDailyLogbookDetail, err := validator.createSchema("update_daily_logbook_detail_schema.json")
	if err != nil {
		return nil, err
	}
	createDailyLogbook, err := validator.createSchema("create_daily_logbook_schema.json")
	if err != nil {
		return nil, err
	}

	validator.RegisterValidator = register
	validator.LoginValidator = login
	validator.MessageValidator = message
	validator.ResendVerificationEmailValidator = resendVerificationEmail
	validator.PasswordResetRequestValidator = passwordResetRequest
	validator.UpdatePasswordValidator = updatePassword
	validator.UpdateEmployeeValidator = updateEmployee
	validator.ChangePasswordValidator = changePassword
	validator.VerifyEmailValidator = verifyEmail
	validator.AddAirlineEmployeeValidator = addAirlineEmployee
	validator.UpdateAirlineEmployeeValidator = updateAirlineEmployee
	validator.CreateLicensePlateValidator = createLicensePlate
	validator.UpdateLicensePlateValidator = updateLicensePlate
	validator.CreateDailyLogbookDetailValidator = createDailyLogbookDetail
	validator.UpdateDailyLogbookDetailValidator = updateDailyLogbookDetail
	validator.CreateDailyLogbookValidator = createDailyLogbook

	updateDailyLogbook, err := validator.createSchema("update_daily_logbook_schema.json")
	if err != nil {
		return nil, err
	}
	validator.UpdateDailyLogbookValidator = updateDailyLogbook

	refreshToken, err := validator.createSchema("refresh_token_schema.json")
	if err != nil {
		return nil, err
	}
	validator.RefreshTokenValidator = refreshToken

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
