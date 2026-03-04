package middleware

import (
	"net/http"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	messagingCache "github.com/champion19/api-flighthours/platform/cache/messaging"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// fallbackInternalErrorMessage is the message returned when no system_messages entry exists
const fallbackInternalErrorMessage = "Internal server error"

var errorToMessageCode = map[error]string{
	domain.ErrDuplicateUser:             domain.MsgUserDuplicate,
	domain.ErrUserCannotSave:            domain.MsgUserCannotSave,
	domain.ErrUserCannotFound:           domain.MsgUserNotFound,
	domain.ErrUserCannotGet:             domain.MsgUserNotFound,
	domain.ErrNotFoundUserByEmail:       domain.MsgUserEmailNotFound,
	domain.ErrGettingUserByEmail:        domain.MsgUserEmailError,
	domain.ErrorEmailNotVerified:        domain.MsgUserEmailNotVerified,
	domain.ErrVerificationTokenNotFound: domain.MsgUserTokenNotFound,
	domain.ErrTokenExpired:              domain.MsgUserTokenExpired,
	domain.ErrTokenAlreadyUsed:          domain.MsgUserTokenUsed,
	domain.ErrRegistrationFailed:        domain.MsgUserRegError,
	domain.ErrRoleRequired:              domain.MsgUserRoleRequired,
	domain.ErrUserCannotDelete:          domain.MsgUserCannotDelete,

	domain.ErrPersonNotFound:     domain.MsgPersonNotFound,
	domain.ErrInvalidTransaction: domain.MsgPersonInvalidTx,

	domain.ErrInvalidJSONFormat:      domain.MsgValJSONInvalid,
	domain.ErrInvalidRequest:         domain.MsgValInvalidReq,
	domain.ErrInvalidID:              domain.MsgValIDInvalid,
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

	domain.ErrRoleAssignmentFailed:       domain.MsgRoleAssignError,
	domain.ErrRoleRemovalFailed:          domain.MsgRoleRemoveError,
	domain.ErrRoleCheckFailed:            domain.MsgRoleCheckError,
	domain.ErrGetUserRolesFailed:         domain.MsgRoleGetError,
	domain.ErrMessageNotFound:            domain.MsgMessageNotFound,
	domain.ErrMessageCodeRequired:        domain.MsgMessageCodeRequired,
	domain.ErrMessageTypeRequired:        domain.MsgMessageTypeRequired,
	domain.ErrMessageTitleRequired:       domain.MsgMessageTitleRequired,
	domain.ErrMessageContentRequired:     domain.MsgMessageContentReq,
	domain.ErrMessageModuleRequired:      domain.MsgMessageModuleRequired,
	domain.ErrMessageCategoryRequired:    domain.MsgMessageCategoryReq,
	domain.ErrMessageCodeDuplicate:       domain.MsgMessageCodeDuplicate,
	domain.ErrMessageCannotSave:          domain.MsgMessageSaveError,
	domain.ErrMessageCannotUpdate:        domain.MsgMessageUpdateError,
	domain.ErrMessageCannotDelete:        domain.MsgMessageDeleteError,
	domain.ErrMessageInvalidType:         domain.MsgMessageInvalidType,
	domain.ErrMessageListFailed:          domain.MsgMessageListError,
	domain.ErrMessageNotRegistered:       domain.MsgMessageNotRegistered,
	domain.ErrMessageInactive:            domain.MsgMessageInactive,
	domain.ErrKeycloakInconsistentState:  domain.MsgKeycloakInconsistentState,
	domain.ErrKeycloakUserCreationFailed: domain.MsgKeycloakCreateError,
	domain.ErrKeycloakCleanupFailed:      domain.MsgKeycloakCleanupError,
	domain.ErrKeycloakUnavailable:        domain.MsgKeycloakUnavailable,
	domain.ErrDatabaseUnavailable:        domain.MsgDatabaseUnavailable,
	domain.ErrIncompleteRegistration:     domain.MsgIncompleteRegistration,
	domain.ErrInvalidToken:               domain.MsgUnauthorized,
	domain.ErrUserNotFound:               domain.MsgUserNotFound,
	domain.ErrDailyLogbookNotFound:       domain.MsgDailyLogbookNotFound,
	domain.ErrDailyLogbookCannotSave:     domain.MsgDailyLogbookSaveError,
	domain.ErrDailyLogbookCannotDelete:   domain.MsgDailyLogbookDeleteError,
	domain.ErrDailyLogbookUnauthorized:   domain.MsgDailyLogbookUnauthorized,
	domain.ErrDailyLogbookInactive:       domain.MsgDailyLogbookInactive,
	domain.ErrTailNumberNotFound:         domain.MsgTailNumberNotFound,
	domain.ErrTailNumberCannotSave:       domain.MsgTailNumberSaveError,
	domain.ErrTailNumberCannotUpdate:     domain.MsgTailNumberUpdateError,
	domain.ErrTailNumberDuplicatePlate:   domain.MsgTailNumberDuplicate,
	domain.ErrTailNumberInvalidModel:     domain.MsgTailNumberInvalidModel,
	domain.ErrTailNumberInvalidAirline:   domain.MsgTailNumberInvalidAirline,
	domain.ErrAirlineRouteNotFound:       domain.MsgAirlineRouteNotFound,
	domain.ErrAirlineRouteCannotSave:     domain.MsgAirlineRouteGetErr,
	domain.ErrAirlineRouteCannotUpdate:   domain.MsgAirlineRouteDeactivateErr,
	domain.ErrAirlineRouteInvalidRoute:   domain.MsgAirlineRouteInvalidRoute,
	domain.ErrAirlineRouteInvalidAirline: domain.MsgAirlineRouteInvalidAirline,
	domain.ErrFlightNotFound:             domain.MsgFlightNotFound,
	domain.ErrFlightCannotSave:           domain.MsgFlightSaveError,
	domain.ErrFlightCannotUpdate:         domain.MsgFlightUpdateError,
	domain.ErrFlightCannotDelete:         domain.MsgFlightDeleteError,
	domain.ErrFlightUnauthorized:         domain.MsgFlightUnauthorized,
	domain.ErrFlightInvalidRoute:         domain.MsgFlightInvalidRoute,
	domain.ErrFlightInvalidLogbook:       domain.MsgFlightInvalidLogbook,
	domain.ErrFlightInvalidTailNumber:    domain.MsgFlightInvalidTailNumber,
	domain.ErrFlightInvalidTimeSequence:  domain.MsgFlightInvalidTimeSequence,
	domain.ErrFlightDuplicate:            domain.MsgFlightDuplicate,
	domain.ErrEngineNotFound:             domain.MsgEngineNotFound,
	domain.ErrManufacturerNotFound:       domain.MsgManufacturerNotFound,
	domain.ErrAircraftModelNotFound:      domain.MsgAircraftModelNotFound,
	domain.ErrRouteNotFound:              domain.MsgRouteNotFound,
	domain.ErrAirlineNotFound:            domain.MsgAirlineNotFound,
	domain.ErrAirportNotFound:            domain.MsgAirportNotFound,
	domain.ErrInternalServer:             domain.MsgServerError,
	domain.ErrRateLimitExceeded:          domain.MsgRateLimitExceeded,
	domain.ErrRefreshTokenFailed:         domain.MsgKCRefreshTokenFailed,
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

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		traceID := GetRequestID(c)
		log := log.WithTraceID(traceID)

		params := buildValidationParams(c)

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
			Message: fallbackInternalErrorMessage,
		})
	}

}

func buildValidationParams(c *gin.Context) []string {
	validationFields, exists := c.Get("validation_fields")
	if !exists {
		return nil
	}

	fields, ok := validationFields.([]string)
	if !ok {
		return nil
	}

	if len(fields) <= 1 {
		return fields
	}

	fieldsStr := fields[0]
	for i := 1; i < len(fields); i++ {
		fieldsStr += ", " + fields[i]
	}
	return []string{fieldsStr}
}
