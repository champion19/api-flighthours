package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// ListDailyLogbooks godoc
// @Summary      List daily logbooks for authenticated employee
// @Description  Returns a list of daily logbooks for the currently authenticated employee
// @Tags         DailyLogbooks
// @Produce      json
// @Param        status query bool false "Filter by status (true for active, false for inactive)"
// @Success      200  {object}  DailyLogbookListResponse
// @Failure      401  {object}  middleware.APIResponse
// @Failure      500  {object}  middleware.APIResponse
// @Router       /daily-logbooks [get]
// @Security     BearerAuth
func (h *handler) ListDailyLogbooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogDailyLogbookList, "action", "list_my_logbooks")

		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogDailyLogbookListError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
			return
		}

		filters := make(map[string]interface{})
		if statusParam := c.Query("status"); statusParam != "" {
			if statusParam == "true" {
				filters["status"] = true
			} else if statusParam == "false" {
				filters["status"] = false
			}
		}

		logbooks, err := h.DailyLogbookInteractor.ListDailyLogbooksByEmployee(c.Request.Context(), employee.ID, filters)
		if err != nil {
			log.Error(logger.LogDailyLogbookListError, "employee_id", employee.ID, "error", err)
			h.Response.Error(c, domain.MsgDailyLogbookListError)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToDailyLogbookListResponse(logbooks, h.EncodeID, baseURL)

		log.Info(logger.LogDailyLogbookListOK, "employee_id", employee.ID, "count", len(logbooks))
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookListOK, response)
	}
}

// GetDailyLogbookByID godoc
// @Summary      Get daily logbook by ID
// @Description  Returns daily logbook information by ID for the authenticated employee (accepts both UUID and obfuscated ID)
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Daily Logbook ID (obfuscated ID)"
// @Success      200  {object}  DailyLogbookResponse
// @Failure      400  {object}  middleware.APIResponse
// @Failure      401  {object}  middleware.APIResponse
// @Failure      403  {object}  middleware.APIResponse
// @Failure      404  {object}  middleware.APIResponse
// @Failure      500  {object}  middleware.APIResponse
// @Router       /daily-logbooks/{id} [get]
// @Security     BearerAuth
func (h *handler) GetDailyLogbookByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookGet, "input_id", inputID)

		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogDailyLogbookGetError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
			return
		}

		// Resolve ID (accepts both UUID and obfuscated ID)
		logbookUUID, responseID := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookGetError, "error", "invalid ID")
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		logbook, err := h.DailyLogbookInteractor.GetDailyLogbookByID(c.Request.Context(), logbookUUID, employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookGetError, "logbook_id", logbookUUID, "error", err)
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
				return
			}
			h.Response.Error(c, domain.MsgDailyLogbookGetErr)
			return
		}

		// Encode employee ID for response
		encodedEmployeeID, _ := h.EncodeID(employee.ID)

		baseURL := GetBaseURL(c)
		response := FromDomainDailyLogbook(logbook, responseID, encodedEmployeeID)
		response.Links = BuildDailyLogbookLinks(baseURL, responseID)

		log.Info(logger.LogDailyLogbookGetOK, "logbook_id", logbookUUID)
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookGetOK, response)
	}
}

// CreateDailyLogbook godoc
// @Summary      Create a new daily logbook
// @Description  Creates a new daily logbook for the authenticated employee
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        request body CreateDailyLogbookRequest true "Daily logbook data"
// @Success      201  {object}  DailyLogbookResponse
// @Failure      400  {object}  middleware.APIResponse
// @Failure      401  {object}  middleware.APIResponse
// @Failure      500  {object}  middleware.APIResponse
// @Router       /daily-logbooks [post]
// @Security     BearerAuth
func (h *handler) CreateDailyLogbook() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogDailyLogbookCreate, "action", "create_logbook")

		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogDailyLogbookCreateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
			return
		}

		var req CreateDailyLogbookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogDailyLogbookCreateError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		// Sanitize input data
		req.Sanitize()

		logbook, err := req.ToDomain(employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookCreateError, "error", err)
			h.Response.Error(c, domain.MsgDailyLogbookSaveError)
			return
		}

		if err := h.DailyLogbookInteractor.CreateDailyLogbook(c.Request.Context(), *logbook); err != nil {
			log.Error(logger.LogDailyLogbookCreateError, "error", err)
			h.Response.Error(c, domain.MsgDailyLogbookSaveError)
			return
		}

		encodedID, err := h.EncodeID(logbook.ID)
		if err != nil {
			h.HandleIDEncodingError(c, logbook.ID, err)
			return
		}

		// Encode employee ID for response
		encodedEmployeeID, _ := h.EncodeID(employee.ID)

		baseURL := GetBaseURL(c)
		response := FromDomainDailyLogbook(logbook, encodedID, encodedEmployeeID)
		response.Links = BuildDailyLogbookCreatedLinks(baseURL, encodedID)

		SetLocationHeader(c, baseURL, "daily-logbooks", encodedID)
		log.Info(logger.LogDailyLogbookCreateOK, "id", logbook.ID)
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookCreated, response)
	}
}
