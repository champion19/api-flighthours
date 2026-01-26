package middleware

import (
	"net/http"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	messagingCache "github.com/champion19/flighthours-api/platform/cache/messaging"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

var errorToMessageCode = map[error]string{
	domain.ErrDuplicateUser:                 domain.MsgUserDuplicate,
	domain.ErrUserCannotSave:                domain.MsgUserCannotSave,
	domain.ErrUserCannotFound:               domain.MsgUserNotFound,
	domain.ErrUserCannotGet:                 domain.MsgUserNotFound,
	domain.ErrNotFoundUserByEmail:           domain.MsgUserEmailNotFound,
	domain.ErrGettingUserByEmail:            domain.MsgUserEmailError,
	domain.ErrorEmailNotVerified:            domain.MsgUserEmailNotVerified,
	domain.ErrVerificationTokenNotFound:     domain.MsgUserTokenNotFound,
	domain.ErrTokenExpired:                  domain.MsgUserTokenExpired,
	domain.ErrTokenAlreadyUsed:              domain.MsgUserTokenUsed,
	domain.ErrRegistrationFailed:            domain.MsgUserRegError,
	domain.ErrRoleRequired:                  domain.MsgUserRoleRequired,
	domain.ErrUserCannotDelete:              domain.MsgUserCannotDelete,

	domain.ErrPersonNotFound:     domain.MsgPersonNotFound,
	domain.ErrInvalidTransaction: domain.MsgPersonInvalidTx,

	domain.ErrInvalidJSONFormat: domain.MsgValJSONInvalid,
	domain.ErrInvalidRequest:    domain.MsgValInvalidReq,
	domain.ErrInvalidID:         domain.MsgValIDInvalid,
	domain.ErrSchemaBadRequest:       domain.MsgValBadFormat,
	domain.ErrSchemaInvalidRequest:   domain.MsgValInvalidReq,
	domain.ErrSchemaReadFailed:       domain.MsgValSchemaRead,
	domain.ErrSchemaEmpty:            domain.MsgValSchemaEmpty,
	domain.ErrSchemaCompileFailed:    domain.MsgValSchemaCompile,
	domain.ErrSchemaValidationFailed: domain.MsgValFailed,
	domain.ErrSchemaBodyReadFailed:   domain.MsgValBodyRead,
	domain.ErrSchemaFieldFormat:      domain.MsgValFieldFormat,
	domain.ErrSchemaFieldRequired:    domain.MsgValFieldRequired,
	domain.ErrSchemaFieldType:        domain.MsgValFieldType,
	domain.ErrSchemaMultipleFields:   domain.MsgValMultiple,

	domain.ErrRoleAssignmentFailed: domain.MsgRoleAssignError,
	domain.ErrRoleRemovalFailed:    domain.MsgRoleRemoveError,
	domain.ErrRoleCheckFailed:      domain.MsgRoleCheckError,
	domain.ErrGetUserRolesFailed:   domain.MsgRoleGetError,
	domain.ErrMessageNotFound:         domain.MsgMessageNotFound,
	domain.ErrMessageCodeRequired:     domain.MsgMessageCodeRequired,
	domain.ErrMessageTypeRequired:     domain.MsgMessageTypeRequired,
	domain.ErrMessageTitleRequired:    domain.MsgMessageTitleRequired,
	domain.ErrMessageContentRequired:  domain.MsgMessageContentReq,
	domain.ErrMessageModuleRequired:   domain.MsgMessageModuleRequired,
	domain.ErrMessageCategoryRequired: domain.MsgMessageCategoryReq,
	domain.ErrMessageCodeDuplicate:    domain.MsgMessageCodeDuplicate,
	domain.ErrMessageCannotSave:       domain.MsgMessageSaveError,
	domain.ErrMessageCannotUpdate:     domain.MsgMessageUpdateError,
	domain.ErrMessageCannotDelete:     domain.MsgMessageDeleteError,
	domain.ErrMessageInvalidType:      domain.MsgMessageInvalidType,
	domain.ErrMessageListFailed:       domain.MsgMessageListError,
	domain.ErrMessageNotRegistered:    domain.MsgMessageNotRegistered,
	domain.ErrMessageInactive:         domain.MsgMessageInactive,
	domain.ErrKeycloakInconsistentState:  domain.MsgKeycloakInconsistentState,
	domain.ErrKeycloakUserCreationFailed: domain.MsgKeycloakCreateError,
	domain.ErrKeycloakCleanupFailed:      domain.MsgKeycloakCleanupError,
	domain.ErrKeycloakUnavailable: domain.MsgKeycloakUnavailable,
	domain.ErrDatabaseUnavailable: domain.MsgDatabaseUnavailable,
	domain.ErrIncompleteRegistration: domain.MsgIncompleteRegistration,
	domain.ErrInvalidToken: domain.MsgUnauthorized,
	domain.ErrUserNotFound: domain.MsgUserNotFound,
	domain.ErrDailyLogbookNotFound:     domain.MsgDailyLogbookNotFound,
	domain.ErrDailyLogbookCannotSave:   domain.MsgDailyLogbookSaveError,
	domain.ErrDailyLogbookCannotUpdate: domain.MsgDailyLogbookUpdateError,
	domain.ErrDailyLogbookCannotDelete: domain.MsgDailyLogbookDeleteError,
	domain.ErrDailyLogbookUnauthorized: domain.MsgDailyLogbookUnauthorized,
	domain.ErrAircraftRegistrationNotFound:       domain.MsgAircraftRegistrationNotFound,
	domain.ErrAircraftRegistrationCannotSave:     domain.MsgAircraftRegistrationSaveError,
	domain.ErrAircraftRegistrationCannotUpdate:   domain.MsgAircraftRegistrationUpdateError,
	domain.ErrAircraftRegistrationDuplicatePlate: domain.MsgAircraftRegistrationDuplicate,
	domain.ErrAircraftRegistrationInvalidModel:   domain.MsgAircraftRegistrationInvalidModel,
	domain.ErrAircraftRegistrationInvalidAirline: domain.MsgAircraftRegistrationInvalidAirline,
	domain.ErrAirlineRouteNotFound:       domain.MsgAirlineRouteNotFound,
	domain.ErrAirlineRouteCannotSave:     domain.MsgAirlineRouteGetErr,
	domain.ErrAirlineRouteCannotUpdate:   domain.MsgAirlineRouteDeactivateErr,
	domain.ErrAirlineRouteInvalidRoute:   domain.MsgAirlineRouteInvalidRoute,
	domain.ErrAirlineRouteInvalidAirline: domain.MsgAirlineRouteInvalidAirline,
	domain.ErrFlightNotFound:            domain.MsgFlightNotFound,
	domain.ErrFlightCannotSave:          domain.MsgFlightSaveError,
	domain.ErrFlightCannotUpdate:        domain.MsgFlightUpdateError,
	domain.ErrFlightCannotDelete:        domain.MsgFlightDeleteError,
	domain.ErrFlightUnauthorized:        domain.MsgFlightUnauthorized,
	domain.ErrFlightInvalidRoute:        domain.MsgFlightInvalidRoute,
	domain.ErrFlightInvalidLogbook:      domain.MsgFlightInvalidLogbook,
	domain.ErrFlightInvalidAircraft:     domain.MsgFlightInvalidAircraft,
	domain.ErrFlightInvalidTimeSequence: domain.MsgFlightInvalidTimeSequence,
	domain.ErrEngineNotFound: domain.MsgEngineNotFound,
	domain.ErrManufacturerNotFound: domain.MsgManufacturerNotFound,
	domain.ErrAircraftModelNotFound: domain.MsgAircraftModelNotFound,
	domain.ErrRouteNotFound: domain.MsgRouteNotFound,
	domain.ErrAirlineNotFound: domain.MsgAirlineNotFound,
	domain.ErrAirportNotFound: domain.MsgAirportNotFound,
	domain.ErrInternalServer: domain.MsgServerError,
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorHandler struct {
	cache *messagingCache.MessageCache
}

var log logger.Logger = logger.NewSlogLogger()

func NewErrorHandler(cache *messagingCache.MessageCache) *ErrorHandler {
	return &ErrorHandler{
		cache: cache,
	}
}

func (h *ErrorHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			traceID := GetRequestID(c)
			log := log.WithTraceID(traceID)

			var params []string
			if validationFields, exists := c.Get("validation_fields"); exists {
				if fields, ok := validationFields.([]string); ok {
					if len(fields) > 1 {
						fieldsStr := fields[0]
						for i := 1; i < len(fields); i++ {
							fieldsStr += ", " + fields[i]
						}
						params = []string{fieldsStr}
					} else {
						params = fields
					}
				}
			}

			if messageCode, ok := errorToMessageCode[err]; ok {
				msg := h.cache.GetMessageResponse(messageCode, params...)
				status := h.cache.GetHTTPStatus(messageCode)

				if msg != nil {
					log.Warn(logger.LogMiddlewareErrorCaught,
						"error", err.Error(),
						"code", msg.Code,
						"status", status,
						"fields", params,
						"path", c.Request.URL.Path,
						"method", c.Request.Method,
						"client_ip", c.ClientIP())

					c.JSON(status, ErrorResponse{
						Success: false,
						Code:    msg.Code,
						Message: msg.Content,
					})
					return
				}
			}

			log.Error(logger.LogMiddlewareInternalErr,
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP())

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Code:    domain.MsgServerError,
				Message: "Error interno del servidor",
			})
		}
	}

}
