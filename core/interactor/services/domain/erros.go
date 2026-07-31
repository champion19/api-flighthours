package domain

import "errors"

var (
	ErrDuplicateUser             = errors.New("Err_DUPLICATE_USER")
	ErrUserCannotSave            = errors.New("Err_USER_CANNOT_SAVE")
	ErrPersonNotFound            = errors.New("Err_PERSON_NOT_FOUND")
	ErrUserNotFound              = errors.New("Err_USER_NOT_FOUND")
	ErrEmailAlreadyVerified      = errors.New("Err_EMAIL_ALREADY_VERIFIED")
	ErrInvalidTransaction        = errors.New("Err_INVALID_TRANSACTION")
	ErrGettingUserByEmail        = errors.New("Err_GETTING_USER_BY_EMAIL")
	ErrNotFoundUserByEmail       = errors.New("Err_NOT_FOUND_USER_BY_EMAIL")
	ErrNotFoundUserById          = errors.New("Err_NOT_FOUND_USER_BY_ID")
	ErrUserCannotFound           = errors.New("Err_USER_CANNOT_FOUND")
	ErrUserCannotGet             = errors.New("Err_USER_CANNOT_GET")
	ErrorEmailNotVerified        = errors.New("Err_EMAIL_NOT_VERIFIED")
	ErrVerificationTokenNotFound = errors.New("Err_VERIFICATION_TOKEN_NOT_FOUND")
	ErrTokenExpired              = errors.New("Err_TOKEN_EXPIRED")
	ErrTokenAlreadyUsed          = errors.New("Err_TOKEN_ALREADY_USED")
	ErrInvalidToken              = errors.New("Err_INVALID_TOKEN")
	ErrRegistrationFailed        = errors.New("Err_REGISTRATION_FAILED")
	ErrRoleRequired              = errors.New("Err_ROLE_REQUIRED")
	ErrDBQueryFailed             = errors.New("Err_DB_QUERY_FAILED")
	ErrUserCannotDelete          = errors.New("Err_USER_CANNOT_DELETE")
	ErrPasswordMismatch          = errors.New("Err_PASSWORD_MISMATCH")
	ErrPasswordUpdateFailed      = errors.New("Err_PASSWORD_UPDATE_FAILED")
	ErrInvalidCurrentPassword    = errors.New("Err_INVALID_CURRENT_PASSWORD")
	ErrUserCannotUpdate          = errors.New("Err_USER_CANNOT_UPDATE")
	ErrKeycloakUpdateFailed      = errors.New("Err_KEYCLOAK_UPDATE_FAILED")
	ErrRoleUpdateFailed          = errors.New("Err_ROLE_UPDATE_FAILED")
	// Data validation errors (from DB constraints)
	ErrInvalidForeignKey = errors.New("Err_INVALID_FOREIGN_KEY") // 1452 - FK constraint fails
	ErrDataTooLong       = errors.New("Err_DATA_TOO_LONG")       // 1406 - Data too long for column
	ErrInvalidData       = errors.New("Err_INVALID_DATA")        // Generic invalid data
	// Date validation errors
	ErrStartDateAfterEndDate = errors.New("Err_START_DATE_AFTER_END_DATE") // start_date > end_date
	ErrInvalidDateFormat     = errors.New("Err_INVALID_DATE_FORMAT")       // Date format is invalid
	ErrEmptyField            = errors.New("Err_EMPTY_FIELD")               // Field contains only whitespace
)

// Infrastructure Errors
var (
	ErrKeycloakInconsistentState  = errors.New("ERR_KC_INCONSISTENT_STATE")
	ErrKeycloakUserCreationFailed = errors.New("ERR_KC_USER_CREATION_FAILED")
	ErrKeycloakCleanupFailed      = errors.New("ERR_KC_CLEANUP_FAILED")
	// Dependency availability errors
	ErrKeycloakUnavailable = errors.New("ERR_KC_UNAVAILABLE")
	ErrDatabaseUnavailable = errors.New("ERR_DB_UNAVAILABLE")
	// Specific user existence errors
	ErrKeycloakUserExists     = errors.New("ERR_KC_USER_EXISTS")
	ErrDatabaseUserExists     = errors.New("ERR_DB_USER_EXISTS")
	ErrIncompleteRegistration = errors.New("ERR_INCOMPLETE_REGISTRATION")
)

// Request Validation errors
var (
	ErrInvalidJSONFormat = errors.New("Err_INVALID_JSON_FORMAT")
	ErrInvalidRequest    = errors.New("Err_INVALID_REQUEST")
	ErrInvalidID         = errors.New("Err_INVALID_ID")
	ErrInternalServer    = errors.New("Err_INTERNAL_SERVER")

	// Schema validation errors
	ErrSchemaBadRequest       = errors.New("ERR_SCHEMA_BAD_REQUEST")
	ErrSchemaInvalidRequest   = errors.New("ERR_SCHEMA_INVALID_REQUEST")
	ErrSchemaReadFailed       = errors.New("ERR_SCHEMA_READ_FAILED")
	ErrSchemaEmpty            = errors.New("ERR_SCHEMA_EMPTY")
	ErrSchemaCompileFailed    = errors.New("ERR_SCHEMA_COMPILE_FAILED")
	ErrSchemaValidationFailed = errors.New("ERR_SCHEMA_VALIDATION_FAILED")
	ErrSchemaBodyReadFailed   = errors.New("ERR_SCHEMA_BODY_READ_FAILED")
	ErrSchemaFieldFormat      = errors.New("ERR_SCHEMA_FIELD_FORMAT")
	ErrSchemaFieldRequired    = errors.New("ERR_SCHEMA_FIELD_REQUIRED")
	ErrSchemaFieldType        = errors.New("ERR_SCHEMA_FIELD_TYPE")
	ErrSchemaMultipleFields   = errors.New("ERR_SCHEMA_MULTIPLE_FIELDS")
	ErrRateLimitExceeded      = errors.New("ERR_RATE_LIMIT_EXCEEDED")
)

// Authorization Errors
var (
	ErrRoleAssignmentFailed = errors.New("Err_ROLE_ASSIGNMENT_FAILED")
	ErrRoleRemovalFailed    = errors.New("Err_ROLE_REMOVAL_FAILED")
	ErrRoleCheckFailed      = errors.New("Err_ROLE_CHECK_FAILED")
	ErrGetUserRolesFailed   = errors.New("Err_GET_USER_ROLES_FAILED")
	ErrRefreshTokenFailed   = errors.New("ERR_REFRESH_TOKEN_FAILED")
)

// Message Management Errors (MOD_M_*)
var (
	ErrMessageNotFound         = errors.New("ERR_MESSAGE_NOT_FOUND")
	ErrMessageCodeRequired     = errors.New("ERR_MESSAGE_CODE_REQUIRED")
	ErrMessageTypeRequired     = errors.New("ERR_MESSAGE_TYPE_REQUIRED")
	ErrMessageTitleRequired    = errors.New("ERR_MESSAGE_TITLE_REQUIRED")
	ErrMessageContentRequired  = errors.New("ERR_MESSAGE_CONTENT_REQUIRED")
	ErrMessageModuleRequired   = errors.New("ERR_MESSAGE_MODULE_REQUIRED")
	ErrMessageCategoryRequired = errors.New("ERR_MESSAGE_CATEGORY_REQUIRED")
	ErrMessageCodeDuplicate    = errors.New("ERR_MESSAGE_CODE_DUPLICATE")
	ErrMessageCannotSave       = errors.New("ERR_MESSAGE_CANNOT_SAVE")
	ErrMessageCannotUpdate     = errors.New("ERR_MESSAGE_CANNOT_UPDATE")
	ErrMessageCannotDelete     = errors.New("ERR_MESSAGE_CANNOT_DELETE")
	ErrMessageInvalidType      = errors.New("ERR_MESSAGE_INVALID_TYPE")
	ErrMessageListFailed       = errors.New("ERR_MESSAGE_LIST_FAILED")
	ErrMessageNotRegistered    = errors.New("ERR_MESSAGE_NOT_REGISTERED")
	ErrMessageInactive         = errors.New("ERR_MESSAGE_INACTIVE")
)

// Airline Management Errors (MOD_AIR_*)
var (
	ErrAirlineNotFound = errors.New("ERR_AIRLINE_NOT_FOUND")
)

// Airport Management Errors (MOD_APT_*)
var (
	ErrAirportNotFound = errors.New("ERR_AIRPORT_NOT_FOUND")
)

// Tail Number Management Errors (MAT_*)
var (
	ErrTailNumberNotFound       = errors.New("ERR_TAIL_NUMBER_NOT_FOUND")
	ErrTailNumberCannotSave     = errors.New("ERR_TAIL_NUMBER_CANNOT_SAVE")
	ErrTailNumberCannotUpdate   = errors.New("ERR_TAIL_NUMBER_CANNOT_UPDATE")
	ErrTailNumberDuplicatePlate = errors.New("ERR_TAIL_NUMBER_DUPLICATE_PLATE")
	ErrTailNumberInvalidModel   = errors.New("ERR_TAIL_NUMBER_INVALID_MODEL")
	ErrTailNumberInvalidAirline = errors.New("ERR_TAIL_NUMBER_INVALID_AIRLINE")
)

// Aircraft Model Management Errors (MOD_AM_*)
var (
	ErrAircraftModelNotFound = errors.New("ERR_AIRCRAFT_MODEL_NOT_FOUND")
)

// Engine Management Errors (MOT_*)
var (
	ErrEngineNotFound = errors.New("ERR_ENGINE_NOT_FOUND")
)

// Manufacturer Management Errors (FAB_*)
var (
	ErrManufacturerNotFound = errors.New("ERR_MANUFACTURER_NOT_FOUND")
)

// Airport Type Management Errors (TAE_*) - Virtual Entity pattern
var (
	ErrAirportTypeNotFound = errors.New("ERR_AIRPORT_TYPE_NOT_FOUND")
)

// Airline Employee Management Errors (EMP_AIR_*)
var (
	ErrAirlineEmployeeNotFound       = errors.New("ERR_AIRLINE_EMPLOYEE_NOT_FOUND")
	ErrAirlineEmployeeCannotSave     = errors.New("ERR_AIRLINE_EMPLOYEE_CANNOT_SAVE")
	ErrAirlineEmployeeCannotUpdate   = errors.New("ERR_AIRLINE_EMPLOYEE_CANNOT_UPDATE")
	ErrAirlineEmployeeInvalidAirline = errors.New("ERR_AIRLINE_EMPLOYEE_INVALID_AIRLINE")
)

// ============================================
// MESSAGE CODES - Constants for use in code
// ============================================

// User Module (MOD_U_*)
const (
	MsgUserDuplicate        = "MOD_U_DUP_ERR_00001"
	MsgUserCannotSave       = "MOD_U_SAVE_ERR_00002"
	MsgUserNotFound         = "MOD_U_GET_ERR_00003"
	MsgUserEmailError       = "MOD_U_EMAIL_ERR_00004"
	MsgUserEmailNotFound    = "MOD_U_EMAIL_NF_ERR_00005"
	MsgUserEmailNotVerified = "MOD_U_EMAIL_NV_ERR_00006"
	MsgUserTokenNotFound    = "MOD_U_TOKEN_NF_ERR_00007"
	MsgUserTokenExpired     = "MOD_U_TOKEN_EXP_ERR_00008"
	MsgUserTokenUsed        = "MOD_U_TOKEN_USED_ERR_00009"
	MsgUserRegError         = "MOD_U_REG_ERR_00010"
	MsgUserRoleRequired     = "MOD_U_ROLE_REQ_ERR_00011"
	MsgUserCannotDelete     = "MOD_U_DEL_ERR_00012"

	MsgUserRegistered    = "MOD_U_REG_EXI_00001"
	MsgUserUpdated       = "MOD_U_UPD_EXI_00002"
	MsgUserDeleted       = "MOD_U_DEL_EXI_00003"
	MsgUserEmailVerified = "MOD_U_VER_EXI_00004"
	MsgUserFound         = "MOD_U_GET_EXI_00005"

	// Update errors
	MsgUserUpdateError         = "MOD_U_UPD_ERR_00013"
	MsgUserKeycloakUpdateError = "MOD_U_KC_UPD_ERR_00014"
	MsgUserRoleUpdateError     = "MOD_U_ROLE_UPD_ERR_00015"
	// Data validation errors (from DB constraints) - These are 400/422, not 500
	MsgInvalidForeignKey = "MOD_V_FK_ERR_00014"   // Invalid reference (e.g., airline doesn't exist)
	MsgDataTooLong       = "MOD_V_LEN_ERR_00015"  // Data exceeds column length
	MsgInvalidData       = "MOD_V_DATA_ERR_00016" // Generic invalid data
)

// Person Module (MOD_P_*)
const (
	MsgPersonNotFound     = "MOD_P_NOT_FOUND_ERR_00001"
	MsgPersonInvalidTx    = "MOD_P_TRANS_ERR_00002"
	MsgPersonRegistered   = "MOD_P_REG_EXI_00001"
	MsgPersonUpdated      = "MOD_P_UPD_EXI_00002"
	MsgPersonCannotDelete = "MOD_P_DEL_ERR_00003"
)

// Validation Module (MOD_V_*)
const (
	MsgValBadFormat     = "MOD_V_VAL_ERR_00001"
	MsgValInvalidReq    = "MOD_V_VAL_ERR_00002"
	MsgValSchemaRead    = "MOD_V_VAL_ERR_00003"
	MsgValSchemaEmpty   = "MOD_V_VAL_ERR_00004"
	MsgValSchemaCompile = "MOD_V_VAL_ERR_00005"
	MsgValFailed        = "MOD_V_VAL_ERR_00006"
	MsgValBodyRead      = "MOD_V_VAL_ERR_00007"
	MsgValFieldFormat   = "MOD_V_VAL_ERR_00008"
	MsgValFieldRequired = "MOD_V_VAL_ERR_00009"
	MsgValFieldType     = "MOD_V_VAL_ERR_00010"
	MsgValMultiple      = "MOD_V_VAL_ERR_00011"
	MsgValJSONInvalid   = "MOD_V_JSON_ERR_00012"
	MsgValIDInvalid     = "MOD_V_ID_ERR_00013"
	// Date validation errors
	MsgValStartDateAfterEndDate = "MOD_V_DATE_ERR_00017"  // start_date > end_date
	MsgValInvalidDateFormat     = "MOD_V_DATE_ERR_00018"  // Invalid date format
	MsgValEmptyField            = "MOD_V_EMPTY_ERR_00019" // Field contains only whitespace
)

// Authorization Module (MOD_A_*)
const (
	MsgRoleAssignError = "MOD_A_ROLE_ASSIGN_ERR_00001"
	MsgRoleRemoveError = "MOD_A_ROLE_REMOVE_ERR_00002"
	MsgRoleCheckError  = "MOD_A_ROLE_CHECK_ERR_00003"
	MsgRoleGetError    = "MOD_A_ROLE_GET_ERR_00004"
	MsgRoleAssigned    = "MOD_A_ROLE_ASSIGN_EXI_00001"
	MsgRoleRemoved     = "MOD_A_ROLE_REMOVE_EXI_00002"
)

// General Module (GEN_*)
const (
	MsgServerError       = "GEN_SRV_ERR_00001"
	MsgUnauthorized      = "GEN_AUTH_ERR_00002"
	MsgForbidden         = "GEN_FORBIDDEN_ERR_00003"
	MsgRateLimitExceeded = "GEN_RATE_LIMIT_ERR_00001"
	MsgOpSuccess         = "GEN_OPE_EXI_00001"
	MsgInfoProcess       = "GEN_INFO_00001"
	MsgWarningAction     = "GEN_WARN_00001"
)

// Message Module (MOD_M_*)
const (
	MsgMessageNotFound       = "MOD_M_NOT_FOUND_ERR_00001"
	MsgMessageCodeRequired   = "MOD_M_CODE_REQ_ERR_00002"
	MsgMessageTypeRequired   = "MOD_M_TYPE_REQ_ERR_00003"
	MsgMessageTitleRequired  = "MOD_M_TITLE_REQ_ERR_00004"
	MsgMessageContentReq     = "MOD_M_CONTENT_REQ_ERR_00005"
	MsgMessageModuleRequired = "MOD_M_MODULE_REQ_ERR_00006"
	MsgMessageCategoryReq    = "MOD_M_CATEGORY_REQ_ERR_00007"
	MsgMessageCodeDuplicate  = "MOD_M_CODE_DUP_ERR_00008"
	MsgMessageSaveError      = "MOD_M_SAVE_ERR_00009"
	MsgMessageUpdateError    = "MOD_M_UPDATE_ERR_00010"
	MsgMessageDeleteError    = "MOD_M_DELETE_ERR_00011"
	MsgMessageInvalidType    = "MOD_M_TYPE_INV_ERR_00012"
	MsgMessageListError      = "MOD_M_LIST_ERR_00013"
	MsgMessageCannotDelete   = "MOD_M_DEL_ERR_00014"
	MsgMessageNotRegistered  = "MOD_M_NOT_REG_ERR_00015"
	MsgMessageInactive       = "MOD_M_INACTIVE_ERR_00016"

	MsgMessageCreated = "MOD_M_CREATE_EXI_00001"
	MsgMessageUpdated = "MOD_M_UPDATE_EXI_00002"
	MsgMessageDeleted = "MOD_M_DELETE_EXI_00003"
	MsgMessageListed  = "MOD_M_LIST_EXI_00004"
)

// Infrastructure Module (MOD_INFRA_*)
const (
	MsgKeycloakInconsistentState = "MOD_INFRA_KC_INCONSISTENT_ERR_00001"
	MsgKeycloakCreateError       = "MOD_INFRA_KC_CREATE_ERR_00002"
	MsgKeycloakCleanupError      = "MOD_INFRA_KC_CLEANUP_ERR_00003"
	// Dependency availability messages
	MsgKeycloakUnavailable = "MOD_INFRA_KC_UNAVAIL_ERR_00004"
	MsgDatabaseUnavailable = "MOD_INFRA_DB_UNAVAIL_ERR_00005"
	MsgDependencyFailure   = "MOD_INFRA_DEP_FAIL_ERR_00006"
	// Specific user existence messages
	MsgKeycloakUserExists     = "MOD_INFRA_KC_USER_EXISTS_ERR_00007"
	MsgDatabaseUserExists     = "MOD_INFRA_DB_USER_EXISTS_ERR_00008"
	MsgIncompleteRegistration = "MOD_INFRA_INCOMPLETE_REG_ERR_00009"
)

// Keycloak Module (MOD_KC_*) - Email Verification and Password Reset
const (
	// Email Verification
	MsgKCEmailVerified        = "MOD_KC_EMAIL_VERIFIED_EXI_00001"
	MsgKCInvalidToken         = "MOD_KC_INVALID_TOKEN_ERR_00001"
	MsgKCEmailVerifyError     = "MOD_KC_EMAIL_VERIFY_ERROR_ERR_00001"
	MsgKCUserNotFound         = "MOD_KC_USER_NOT_FOUND_ERR_00001"
	MsgKCEmailAlreadyVerified = "MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001"
	// Verification Email Sending
	MsgKCVerifEmailSent   = "MOD_KC_VERIF_EMAIL_SENT_EXI_00001"
	MsgKCVerifEmailError  = "MOD_KC_VERIF_EMAIL_ERROR_ERR_00001"
	MsgKCVerifEmailResent = "MOD_KC_VERIF_EMAIL_RESENT_EXI_00001"
	// Password Reset
	MsgKCPwdResetSent  = "MOD_KC_PWD_RESET_SENT_EXI_00001"
	MsgKCPwdResetError = "MOD_KC_PWD_RESET_ERROR_ERR_00001"
	// Password Update
	MsgKCPwdUpdated            = "MOD_KC_PWD_UPDATED_EXI_00001"
	MsgKCPwdUpdateError        = "MOD_KC_PWD_UPDATE_ERROR_ERR_00001"
	MsgKCPwdMismatch           = "MOD_KC_PWD_MISMATCH_ERR_00001"
	MsgKCPwdUpdateTokenInvalid = "MOD_KC_PWD_UPDATE_TOKEN_INVALID_ERR_00001"
	// Change Password (authenticated user changes their own password)
	MsgKCPwdChanged           = "MOD_KC_PWD_CHANGED_EXI_00001"
	MsgKCPwdChangeError       = "MOD_KC_PWD_CHANGE_ERROR_ERR_00001"
	MsgKCPwdCurrentInvalid    = "MOD_KC_PWD_CURRENT_INVALID_ERR_00001"
	MsgKCPwdChangeNewMismatch = "MOD_KC_PWD_CHANGE_MISMATCH_ERR_00001"
	// Login
	MsgKCLoginEmailNotVerified = "MOD_KC_LOGIN_EMAIL_NOT_VERIFIED_ERR_00001"
	MsgKCLoginSuccess          = "MOD_KC_LOGIN_SUCCESS_EXI_00001"
	// Refresh Token
	MsgKCRefreshTokenFailed  = "MOD_KC_REFRESH_TOKEN_ERR_00001"
	MsgKCRefreshTokenSuccess = "MOD_KC_REFRESH_TOKEN_EXI_00001"
)

// Airline Module (MOD_AIR_*)
const (
	// Errors
	MsgAirlineNotFound      = "MOD_AIR_NOT_FOUND_ERR_00001"
	MsgAirlineActivateErr   = "MOD_AIR_ACTIVATE_ERR_00002"
	MsgAirlineDeactivateErr = "MOD_AIR_DEACTIVATE_ERR_00003"
	MsgAirlineListError     = "MOD_AIR_LIST_ERR_00004"
	// Success
	MsgAirlineGetOK        = "MOD_AIR_GET_EXI_00001"
	MsgAirlineActivateOK   = "MOD_AIR_ACTIVATE_EXI_00002"
	MsgAirlineDeactivateOK = "MOD_AIR_DEACTIVATE_EXI_00003"
	MsgAirlineListOK       = "MOD_AIR_LIST_EXI_00004"
)

// Airport Module (MOD_APT_*)
const (
	// Errors
	MsgAirportNotFound      = "MOD_APT_NOT_FOUND_ERR_00001"
	MsgAirportActivateErr   = "MOD_APT_ACTIVATE_ERR_00002"
	MsgAirportDeactivateErr = "MOD_APT_DEACTIVATE_ERR_00003"
	// Success
	MsgAirportGetOK        = "MOD_APT_GET_EXI_00001"
	MsgAirportActivateOK   = "MOD_APT_ACTIVATE_EXI_00002"
	MsgAirportDeactivateOK = "MOD_APT_DEACTIVATE_EXI_00003"
)

// City Module (CIU_*) - Ciudad ( - Virtual Entity pattern)
const (
	// ========================================
	// Consultar - CIU_CON_*
	// ========================================
	MsgCityGetOK    = "CIU_CON_EXI_01301" // Éxito - Ciudad consultada (aeropuertos en la ciudad)
	MsgCityNotFound = "CIU_CON_ERR_01302" // Error - Ciudad no encontrada (sin aeropuertos)
	MsgCityGetErr   = "CIU_CON_ERR_01303" // Error - Error técnico al consultar
)

// Country Module (PAI_*) - País - Virtual Entity pattern
const (
	// ========================================
	// Consultar  - PAI_CON_*
	// ========================================
	MsgCountryGetOK    = "PAI_CON_EXI_03801" // Éxito - País consultado (aeropuertos en el país)
	MsgCountryNotFound = "PAI_CON_ERR_03802" // Error - País no encontrado (sin aeropuertos)
	MsgCountryGetErr   = "PAI_CON_ERR_03803" // Error - Error técnico al consultar
)

// Airport Type Module (TAE_*) - Tipo Aeropuerto (- Virtual Entity pattern)
const (
	// ========================================
	// Consultar  - TAE_CON_*
	// ========================================
	MsgAirportTypeGetOK    = "TAE_CON_EXI_04601" // Éxito - Tipo de aeropuerto consultado (aeropuertos de ese tipo)
	MsgAirportTypeNotFound = "TAE_CON_ERR_04602" // Error - Tipo de aeropuerto no encontrado (sin aeropuertos)
	MsgAirportTypeGetErr   = "TAE_CON_ERR_04603" // Error - Error técnico al consultar
)

// Crew Member Type Module (TIN_*) - Tipo Integrante
const (
	// ========================================
	// Consultar  - TIN_CON_*
	// ========================================
	MsgCrewMemberTypeGetOK    = "TIN_CON_EXI_04701" // Éxito - Tipo de integrante consultado
	MsgCrewMemberTypeNotFound = "TIN_CON_ERR_04702" // Error - Tipo de integrante no encontrado (sin empleados)
	MsgCrewMemberTypeGetErr   = "TIN_CON_ERR_04703" // Error - Error técnico al consultar
)

// Crew Member Module (TRI_*) - Tripulante (roster de personas por piloto: Primer Oficial, tripulación de cabina)
const (
	// ========================================
	// Agregar  - TRI_AGR_*
	// ========================================
	MsgCrewMemberCreated = "TRI_AGR_EXI_05101" // Éxito - Tripulante creado o ya existente en el roster

	// ========================================
	// Consultar  - TRI_CON_*
	// ========================================
	MsgCrewMemberSearchOK = "TRI_CON_EXI_05201" // Éxito - Búsqueda de tripulantes exitosa
)

// Tail Number Module (MAT_*) - Matrícula
const (
	// ========================================
	// Consultar - MAT_CON_*
	// ========================================
	MsgTailNumberGetOK    = "MAT_CON_EXI_03301" // Éxito - Matrícula consultada
	MsgTailNumberNotFound = "MAT_CON_ERR_03302" // Error - Matrícula no encontrada
	MsgTailNumberGetErr   = "MAT_CON_ERR_03303" // Error - Error técnico al consultar

	// ========================================
	// Agregar  - MAT_AGR_*
	// ========================================
	MsgTailNumberCreated   = "MAT_AGR_EXI_03401" // Éxito - Matrícula creada
	MsgTailNumberSaveError = "MAT_AGR_ERR_03402" // Error - Error técnico al crear
	MsgTailNumberDuplicate = "MAT_AGR_ERR_03403" // Error - Matrícula duplicada

	// ========================================
	// Editar  - MAT_EDI_*
	// ========================================
	MsgTailNumberUpdated     = "MAT_EDI_EXI_03501" // Éxito - Matrícula actualizada
	MsgTailNumberUpdateError = "MAT_EDI_ERR_03502" // Error - Error técnico al editar

	// ========================================
	// Listar - MAT_LIST_*
	// ========================================
	MsgTailNumberListOK    = "MAT_LIST_EXI_03001" // Éxito - Lista de matrículas obtenida
	MsgTailNumberListError = "MAT_LIST_ERR_03002" // Error - Error al listar matrículas

	// ========================================
	// Validaciones - MAT_VAL_*
	// ========================================
	MsgTailNumberInvalidModel   = "MAT_VAL_ERR_03601" // Error - Modelo de aeronave inválido
	MsgTailNumberInvalidAirline = "MAT_VAL_ERR_03602" // Error - Aerolínea inválida
)

// Aircraft Model Module (MOD_AM_*) - Modelo de Aeronave
const (
	// ========================================
	// Consultar (HU36) - MOD_AM_CON_*
	// ========================================
	MsgAircraftModelGetOK    = "MOD_AM_CON_EXI_03601" // Éxito - Modelo de aeronave consultado
	MsgAircraftModelNotFound = "MOD_AM_CON_ERR_03602" // Error - Modelo de aeronave no encontrado
	MsgAircraftModelGetErr   = "MOD_AM_CON_ERR_03603" // Error - Error técnico al consultar

	// ========================================
	// Listar tipos  - MOD_AM_LIST_*
	// ========================================
	MsgAircraftModelListOK    = "MOD_AM_LIST_EXI_04301" // Éxito - Lista de modelos/tipos obtenida
	MsgAircraftModelListError = "MOD_AM_LIST_ERR_04302" // Error - Error al listar modelos/tipos

	// ========================================
	// Inactivar  - MOD_AM_INA_*
	// ========================================
	MsgAircraftModelDeactivateOK  = "MOD_AM_INA_EXI_04101" // Éxito - Modelo de aeronave inactivado
	MsgAircraftModelDeactivateErr = "MOD_AM_INA_ERR_04102" // Error - Error técnico al inactivar

	// ========================================
	// Activar  - MOD_AM_ACT_*
	// ========================================
	MsgAircraftModelActivateOK  = "MOD_AM_ACT_EXI_04201" // Éxito - Modelo de aeronave activado
	MsgAircraftModelActivateErr = "MOD_AM_ACT_ERR_04202" // Error - Error técnico al activar
)

// Aircraft Family Module (FAM_*) - Familia de Aeronave (HU32)
const (
	// ========================================
	// Consultar - FAM_CON_*
	// ========================================
	MsgAircraftFamilyGetOK    = "FAM_CON_EXI_03201" // Éxito - Familia de aeronaves consultada
	MsgAircraftFamilyNotFound = "FAM_CON_ERR_03202" // Error - Familia de aeronaves no encontrada
	MsgAircraftFamilyGetErr   = "FAM_CON_ERR_03203" // Error - Error técnico al consultar familia
)

// Engine Module (MOT_*) - Motor
const (
	// ========================================
	// Consultar  - MOT_CON_*
	// ========================================
	MsgEngineGetOK    = "MOT_CON_EXI_03701" // Éxito - Motor consultado exitosamente
	MsgEngineNotFound = "MOT_CON_ERR_03702" // Error - Motor no encontrado
	MsgEngineGetErr   = "MOT_CON_ERR_03703" // Error - Error técnico al consultar

	// ========================================
	// Listar - MOT_LIST_*
	// ========================================
	MsgEngineListOK    = "MOT_LIST_EXI_03704" // Éxito - Lista de motores obtenida
	MsgEngineListError = "MOT_LIST_ERR_03705" // Error - Error al listar motores
)

// Manufacturer Module (FAB_*) - Fabricante
const (
	// ========================================
	// Consultar (HU31) - FAB_CON_*
	// ========================================
	MsgManufacturerGetOK    = "FAB_CON_EXI_03101" // Éxito - Fabricante consultado exitosamente
	MsgManufacturerNotFound = "FAB_CON_ERR_03102" // Error - Fabricante no encontrado
	MsgManufacturerGetErr   = "FAB_CON_ERR_03103" // Error - Error técnico al consultar

	// ========================================
	// Listar - FAB_LIST_*
	// ========================================
	MsgManufacturerListOK    = "FAB_LIST_EXI_03104" // Éxito - Lista de fabricantes obtenida
	MsgManufacturerListError = "FAB_LIST_ERR_03105" // Error - Error al listar fabricantes
)

// DailyLogbook Management Errors (BIT_*)
var (
	ErrDailyLogbookNotFound     = errors.New("ERR_DAILY_LOGBOOK_NOT_FOUND")
	ErrDailyLogbookCannotSave   = errors.New("ERR_DAILY_LOGBOOK_CANNOT_SAVE")
	ErrDailyLogbookCannotUpdate = errors.New("ERR_DAILY_LOGBOOK_CANNOT_UPDATE")
	ErrDailyLogbookCannotDelete = errors.New("ERR_DAILY_LOGBOOK_CANNOT_DELETE")
	ErrDailyLogbookUnauthorized = errors.New("ERR_DAILY_LOGBOOK_UNAUTHORIZED")
	ErrDailyLogbookInactive     = errors.New("ERR_DAILY_LOGBOOK_INACTIVE")
)

// DailyLogbook Module (BIT_*) - Bitácora Diaria
const (
	// ========================================
	// Consultar  - BIT_CON_*
	// ========================================
	MsgDailyLogbookGetOK    = "BIT_CON_EXI_01901" // Éxito - Bitácora consultada
	MsgDailyLogbookNotFound = "BIT_CON_ERR_01903" // Error - Bitácora no encontrada
	MsgDailyLogbookGetErr   = "BIT_CON_ERR_01904" // Error - Error técnico al consultar

	// ========================================
	// Agregar - BIT_AGR_*
	// ========================================
	MsgDailyLogbookCreated   = "BIT_AGR_EXI_01801" // Éxito - Bitácora creada
	MsgDailyLogbookSaveError = "BIT_AGR_ERR_01804" // Error - Error técnico al crear

	// ========================================
	// Editar - BIT_EDI_*
	// ========================================
	MsgDailyLogbookUpdated     = "BIT_EDI_EXI_01701" // Éxito - Bitácora actualizada
	MsgDailyLogbookUpdateError = "BIT_EDI_ERR_01704" // Error - Error técnico al editar

	// ========================================
	// Eliminar - BIT_DEL_*
	// ========================================
	MsgDailyLogbookDeleted     = "BIT_DEL_EXI_01601" // Éxito - Bitácora eliminada
	MsgDailyLogbookDeleteError = "BIT_DEL_ERR_01604" // Error - Error técnico al eliminar

	// ========================================
	// Activar - BIT_ACT_*
	// ========================================
	MsgDailyLogbookActivateOK  = "BIT_ACT_EXI_01501" // Éxito - Bitácora activada
	MsgDailyLogbookActivateErr = "BIT_ACT_ERR_01504" // Error - Error técnico al activar

	// ========================================
	// Inactivar - BIT_INA_*
	// ========================================
	MsgDailyLogbookDeactivateOK  = "BIT_INA_EXI_01401" // Éxito - Bitácora inactivada
	MsgDailyLogbookDeactivateErr = "BIT_INA_ERR_01404" // Error - Error técnico al inactivar

	// ========================================
	// Listar - BIT_LIST_*
	// ========================================
	MsgDailyLogbookListOK    = "BIT_LIST_EXI_01001" // Éxito - Lista de bitácoras obtenida
	MsgDailyLogbookListError = "BIT_LIST_ERR_01002" // Error - Error al listar bitácoras

	// ========================================
	// Autorización - BIT_AUTH_*
	// ========================================
	MsgDailyLogbookUnauthorized = "BIT_AUTH_ERR_00001" // Error - No autorizado para esta bitácora
	MsgDailyLogbookInactive     = "BIT_VAL_ERR_01905"  // Error - Bitácora cerrada, no se pueden agregar ni editar vuelos
)

// Flight Management Errors (VUE_*) - Also used for DailyLogbookDetail
var (
	ErrFlightNotFound            = errors.New("ERR_FLIGHT_NOT_FOUND")
	ErrFlightCannotSave          = errors.New("ERR_FLIGHT_CANNOT_SAVE")
	ErrFlightCannotUpdate        = errors.New("ERR_FLIGHT_CANNOT_UPDATE")
	ErrFlightCannotDelete        = errors.New("ERR_FLIGHT_CANNOT_DELETE")
	ErrFlightUnauthorized        = errors.New("ERR_FLIGHT_UNAUTHORIZED")
	ErrFlightInvalidAirport      = errors.New("ERR_FLIGHT_INVALID_AIRPORT")
	ErrFlightInvalidLogbook      = errors.New("ERR_FLIGHT_INVALID_LOGBOOK")
	ErrFlightInvalidTailNumber   = errors.New("ERR_FLIGHT_INVALID_TAIL_NUMBER")
	ErrFlightInvalidTimeSequence = errors.New("ERR_FLIGHT_INVALID_TIME_SEQUENCE")
	ErrFlightDuplicate           = errors.New("ERR_FLIGHT_DUPLICATE")
	ErrFlightInvalidApproach     = errors.New("ERR_FLIGHT_INVALID_APPROACH")
)

// Flight Module (VUE_*) - Vuelo
const (
	// ========================================
	// Consultar - VUE_CON_*
	// ========================================
	MsgFlightGetOK    = "VUE_CON_EXI_04801" // Éxito - Vuelo consultado exitosamente
	MsgFlightNotFound = "VUE_CON_ERR_04803" // Error - Vuelo no encontrado
	MsgFlightGetErr   = "VUE_CON_ERR_04804" // Error - Error técnico al consultar

	// ========================================
	// Editar - VUE_EDI_*
	// ========================================
	MsgFlightUpdated     = "VUE_EDI_EXI_04901" // Éxito - Vuelo actualizado
	MsgFlightUpdateError = "VUE_EDI_ERR_04904" // Error - Error técnico al editar

	// ========================================
	// Registrar - VUE_REG_*
	// ========================================
	MsgFlightCreated   = "VUE_REG_EXI_05001" // Éxito - Vuelo registrado
	MsgFlightSaveError = "VUE_REG_ERR_05004" // Error - Error técnico al registrar
	MsgFlightDuplicate = "VUE_REG_ERR_05005" // Error - Vuelo duplicado

	// ========================================
	// Listar - VUE_LIST_*
	// ========================================
	MsgFlightListOK    = "VUE_LIST_EXI_04800" // Éxito - Lista de vuelos obtenida
	MsgFlightListError = "VUE_LIST_ERR_04802" // Error - Error al listar vuelos

	// ========================================
	// Validaciones - VUE_VAL_*
	// ========================================
	MsgFlightInvalidAirport      = "VUE_VAL_ERR_04805" // Error - Aeropuerto de origen/destino inválido
	MsgFlightInvalidLogbook      = "VUE_VAL_ERR_04806" // Error - Bitácora inválida
	MsgFlightInvalidTailNumber   = "VUE_VAL_ERR_04807" // Error - Matrícula de aeronave inválida
	MsgFlightInvalidTimeSequence = "VUE_VAL_ERR_04808" // Error - Secuencia de tiempos inválida (out < takeoff < landing < in)
	MsgFlightInvalidApproach     = "VUE_VAL_ERR_04809" // Error - Combinación de categoría/subtipo/autoland de aproximación inválida

	// ========================================
	// Eliminar - VUE_DEL_*
	// ========================================
	MsgFlightDeleted     = "VUE_DEL_EXI_01801" // Éxito - Vuelo eliminado exitosamente
	MsgFlightDeleteError = "VUE_DEL_ERR_01804" // Error - Error técnico al eliminar

	// ========================================
	// Autorización - VUE_AUTH_*
	// ========================================
	MsgFlightUnauthorized = "VUE_AUTH_ERR_00001" // Error - No autorizado para este vuelo
)

// Airline Employee Module (EMP_AIR_*) - Empleado Aerolínea (Release 15)
const (
	// ========================================
	// Consultar - EMP_AIR_CON_*
	// ========================================
	MsgAirlineEmployeeGetOK    = "EMP_AIR_CON_EXI_02601" // Éxito - Empleado aerolínea consultado
	MsgAirlineEmployeeNotFound = "EMP_AIR_CON_ERR_02602" // Error - Empleado aerolínea no encontrado
	MsgAirlineEmployeeGetErr   = "EMP_AIR_CON_ERR_02603" // Error - Error técnico al consultar

	// ========================================
	// Editar - EMP_AIR_EDI_*
	// ========================================
	MsgAirlineEmployeeUpdated     = "EMP_AIR_EDI_EXI_02701" // Éxito - Empleado aerolínea actualizado
	MsgAirlineEmployeeUpdateError = "EMP_AIR_EDI_ERR_02702" // Error - Error técnico al editar

	// ========================================
	// Agregar - EMP_AIR_AGR_*
	// ========================================
	MsgAirlineEmployeeCreated   = "EMP_AIR_AGR_EXI_02801" // Éxito - Empleado aerolínea creado
	MsgAirlineEmployeeSaveError = "EMP_AIR_AGR_ERR_02802" // Error - Error técnico al crear
	MsgAirlineEmployeeDuplicate = "EMP_AIR_AGR_ERR_02803" // Error - Empleado duplicado

	// ========================================
	// Activar - EMP_AIR_ACT_*
	// ========================================
	MsgAirlineEmployeeActivateOK  = "EMP_AIR_ACT_EXI_02901" // Éxito - Empleado aerolínea activado
	MsgAirlineEmployeeActivateErr = "EMP_AIR_ACT_ERR_02902" // Error - Error técnico al activar

	// ========================================
	// Inactivar - EMP_AIR_INA_*
	// ========================================
	MsgAirlineEmployeeDeactivateOK  = "EMP_AIR_INA_EXI_03001" // Éxito - Empleado aerolínea inactivado
	MsgAirlineEmployeeDeactivateErr = "EMP_AIR_INA_ERR_03002" // Error - Error técnico al inactivar

	// ========================================
	// Listar - EMP_AIR_LIST_*
	// ========================================
	MsgAirlineEmployeeListOK    = "EMP_AIR_LIST_EXI_02601" // Éxito - Lista de empleados aerolínea obtenida
	MsgAirlineEmployeeListError = "EMP_AIR_LIST_ERR_02602" // Error - Error al listar empleados aerolínea

	// ========================================
	// Validaciones - EMP_AIR_VAL_*
	// ========================================
	MsgAirlineEmployeeInvalidAirline = "EMP_AIR_VAL_ERR_02604" // Error - Aerolínea inválida
)

// Flight Summary Module (HVU_*) - Horas de Vuelo
const (
	// ========================================
	// Consultar - HVU_CON_*
	// ========================================
	MsgFlightSummaryGetOK   = "HVU_CON_EXI_05101" // Éxito - Resumen de horas consultado
	MsgFlightSummaryGetErr  = "HVU_CON_ERR_05102" // Error - Error técnico al consultar resumen
	MsgFlightSummaryInvalid = "HVU_CON_ERR_05103" // Error - Período o fecha inválidos

	// ========================================
	// Alertas - HVU_ALR_*
	// ========================================
	MsgFlightAlertsGetOK  = "HVU_ALR_EXI_05201" // Éxito - Alertas consultadas
	MsgFlightAlertsGetErr = "HVU_ALR_ERR_05202" // Error - Error técnico al consultar alertas

	// ========================================
	// Vuelos Recientes - HVU_REC_*
	// ========================================
	MsgRecentFlightsGetOK  = "HVU_REC_EXI_05301" // Éxito - Vuelos recientes consultados
	MsgRecentFlightsGetErr = "HVU_REC_ERR_05302" // Error - Error técnico al consultar vuelos recientes
)
