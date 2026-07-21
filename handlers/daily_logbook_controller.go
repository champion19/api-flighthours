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
		log := log.WithTraceID(traceID)

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
		log := log.WithTraceID(traceID)
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
			log.Warn(logger.LogDailyLogbookGetError, "error", errInvalidID)
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
		h.encodeTailNumberID(&response)
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
		log := log.WithTraceID(traceID)

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

		// Resolve tail_number_id (accepts both UUID and obfuscated ID) before persisting
		if req.TailNumberID != nil && *req.TailNumberID != "" {
			tailNumberUUID, _ := h.resolveID(*req.TailNumberID)
			if tailNumberUUID == "" {
				log.Warn(logger.LogDailyLogbookCreateError, "error", "invalid tail number ID")
				h.Response.Error(c, domain.MsgFlightInvalidTailNumber)
				return
			}
			req.TailNumberID = &tailNumberUUID
		}

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
		h.encodeTailNumberID(&response)
		response.Links = BuildDailyLogbookCreatedLinks(baseURL, encodedID)

		SetLocationHeader(c, baseURL, "daily-logbooks", encodedID)
		log.Info(logger.LogDailyLogbookCreateOK, "id", logbook.ID)
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookCreated, response)
	}
}

// UpdateDailyLogbook godoc
// @Summary      Update a daily logbook
// @Description  Updates an existing daily logbook for the authenticated employee (accepts both UUID and obfuscated ID)
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Daily Logbook ID (obfuscated ID)"
// @Param        request body UpdateDailyLogbookRequest true "Updated daily logbook data"
// @Success      200  {object}  DailyLogbookResponse
// @Failure      400  {object}  middleware.APIResponse
// @Failure      401  {object}  middleware.APIResponse
// @Failure      403  {object}  middleware.APIResponse
// @Failure      404  {object}  middleware.APIResponse
// @Failure      500  {object}  middleware.APIResponse
// @Router       /daily-logbooks/{id} [put]
// @Security     BearerAuth
func (h *handler) UpdateDailyLogbook() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogDailyLogbookUpdateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
			return
		}

		inputID := c.Param("id")
		log.Info(logger.LogDailyLogbookUpdate, "input_id", inputID)

		// Resolve ID (accepts both UUID and obfuscated ID)
		logbookUUID, responseID := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookUpdateError, "error", errInvalidID)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		var req UpdateDailyLogbookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogDailyLogbookUpdateError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		// Sanitize input data
		req.Sanitize()

		// Resolve tail_number_id (accepts both UUID and obfuscated ID) before persisting
		if req.TailNumberID != nil && *req.TailNumberID != "" {
			tailNumberUUID, _ := h.resolveID(*req.TailNumberID)
			if tailNumberUUID == "" {
				log.Warn(logger.LogDailyLogbookUpdateError, "error", "invalid tail number ID")
				h.Response.Error(c, domain.MsgFlightInvalidTailNumber)
				return
			}
			req.TailNumberID = &tailNumberUUID
		}

		logbook, err := req.ToDomain(logbookUUID, employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookUpdateError, "error", err)
			h.Response.Error(c, domain.MsgDailyLogbookUpdateError)
			return
		}

		if err := h.DailyLogbookInteractor.UpdateDailyLogbook(c.Request.Context(), *logbook, employee.ID); err != nil {
			log.Error(logger.LogDailyLogbookUpdateError, "error", err)
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
				return
			}
			h.Response.Error(c, domain.MsgDailyLogbookUpdateError)
			return
		}

		// Encode employee ID for response
		encodedEmployeeID, _ := h.EncodeID(employee.ID)

		baseURL := GetBaseURL(c)
		response := FromDomainDailyLogbook(logbook, responseID, encodedEmployeeID)
		h.encodeTailNumberID(&response)
		response.Links = BuildDailyLogbookLinks(baseURL, responseID)

		log.Info(logger.LogDailyLogbookUpdateOK, "logbook_id", logbookUUID)
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookUpdated, response)
	}
}

// DeleteDailyLogbook godoc
// @Summary      Delete a daily logbook
// @Description  Deletes an existing daily logbook for the authenticated employee (accepts both UUID and obfuscated ID)
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Daily Logbook ID (obfuscated ID)"
// @Success      200  {object}  DailyLogbookDeleteResponse
// @Failure      400  {object}  middleware.APIResponse
// @Failure      401  {object}  middleware.APIResponse
// @Failure      403  {object}  middleware.APIResponse
// @Failure      404  {object}  middleware.APIResponse
// @Failure      500  {object}  middleware.APIResponse
// @Router       /daily-logbooks/{id} [delete]
// @Security     BearerAuth
func (h *handler) DeleteDailyLogbook() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogDailyLogbookDeleteError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
			return
		}

		inputID := c.Param("id")
		log.Info(logger.LogDailyLogbookDelete, "input_id", inputID)

		// Resolve ID (accepts both UUID and obfuscated ID)
		logbookUUID, responseID := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookDeleteError, "error", errInvalidID)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Verify logbook exists and belongs to employee
		_, err := h.DailyLogbookInteractor.GetDailyLogbookByID(c.Request.Context(), logbookUUID, employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookGetError, "logbook_id", logbookUUID, "error", err)
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
				return
			}
			h.Response.Error(c, domain.MsgDailyLogbookGetErr)
			return
		}

		if err := h.DailyLogbookInteractor.DeleteDailyLogbook(c.Request.Context(), logbookUUID); err != nil {
			log.Error(logger.LogDailyLogbookDeleteError, "logbook_id", logbookUUID, "error", err)
			h.Response.Error(c, domain.MsgDailyLogbookDeleteError)
			return
		}

		baseURL := GetBaseURL(c)
		response := DailyLogbookDeleteResponse{
			ID:      responseID,
			Deleted: true,
			Links:   BuildDailyLogbookDeletedLinks(baseURL),
		}

		log.Info(logger.LogDailyLogbookDeleteOK, "logbook_id", logbookUUID)
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookDeleted, response)
	}
}

// ActivateDailyLogbook godoc
// @Summary      Activate a daily logbook
// @Description  Sets daily logbook status to active (accepts both UUID and obfuscated ID)
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Daily Logbook ID (obfuscated ID)"
// @Success      200  {object}  DailyLogbookStatusResponse
// @Failure      400  {object}  middleware.APIResponse
// @Failure      401  {object}  middleware.APIResponse
// @Failure      403  {object}  middleware.APIResponse
// @Failure      404  {object}  middleware.APIResponse
// @Failure      500  {object}  middleware.APIResponse
// @Router       /daily-logbooks/{id}/activate [patch]
// @Security     BearerAuth
func (h *handler) ActivateDailyLogbook() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogDailyLogbookActivateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
			return
		}

		inputID := c.Param("id")
		log.Info(logger.LogDailyLogbookActivate, "input_id", inputID)

		// Resolve ID (accepts both UUID and obfuscated ID)
		logbookUUID, responseID := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookActivateError, "error", errInvalidID)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Verify logbook exists and belongs to employee
		_, err := h.DailyLogbookInteractor.GetDailyLogbookByID(c.Request.Context(), logbookUUID, employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookGetError, "logbook_id", logbookUUID, "error", err)
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
				return
			}
			h.Response.Error(c, domain.MsgDailyLogbookGetErr)
			return
		}

		if err := h.DailyLogbookInteractor.ActivateDailyLogbook(c.Request.Context(), logbookUUID); err != nil {
			log.Error(logger.LogDailyLogbookActivateError, "logbook_id", logbookUUID, "error", err)
			h.Response.Error(c, domain.MsgDailyLogbookActivateErr)
			return
		}

		baseURL := GetBaseURL(c)
		response := DailyLogbookStatusResponse{
			ID:      responseID,
			Status:  "active",
			Updated: true,
			Links:   BuildDailyLogbookStatusLinks(baseURL, responseID, true),
		}

		log.Info(logger.LogDailyLogbookActivateOK, "logbook_id", logbookUUID)
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookActivateOK, response)
	}
}

// DeactivateDailyLogbook godoc
// @Summary      Deactivate a daily logbook
// @Description  Sets daily logbook status to inactive (accepts both UUID and obfuscated ID)
// @Tags         DailyLogbooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Daily Logbook ID (obfuscated ID)"
// @Success      200  {object}  DailyLogbookStatusResponse
// @Failure      400  {object}  middleware.APIResponse
// @Failure      401  {object}  middleware.APIResponse
// @Failure      403  {object}  middleware.APIResponse
// @Failure      404  {object}  middleware.APIResponse
// @Failure      500  {object}  middleware.APIResponse
// @Router       /daily-logbooks/{id}/deactivate [patch]
// @Security     BearerAuth
func (h *handler) DeactivateDailyLogbook() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		// Get authenticated employee from context
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogDailyLogbookDeactivateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
			return
		}

		inputID := c.Param("id")
		log.Info(logger.LogDailyLogbookDeactivate, "input_id", inputID)

		// Resolve ID (accepts both UUID and obfuscated ID)
		logbookUUID, responseID := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookDeactivateError, "error", errInvalidID)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Verify logbook exists and belongs to employee
		_, err := h.DailyLogbookInteractor.GetDailyLogbookByID(c.Request.Context(), logbookUUID, employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookGetError, "logbook_id", logbookUUID, "error", err)
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgDailyLogbookUnauthorized)
				return
			}
			h.Response.Error(c, domain.MsgDailyLogbookGetErr)
			return
		}

		if err := h.DailyLogbookInteractor.DeactivateDailyLogbook(c.Request.Context(), logbookUUID); err != nil {
			log.Error(logger.LogDailyLogbookDeactivateError, "logbook_id", logbookUUID, "error", err)
			h.Response.Error(c, domain.MsgDailyLogbookDeactivateErr)
			return
		}

		baseURL := GetBaseURL(c)
		response := DailyLogbookStatusResponse{
			ID:      responseID,
			Status:  "inactive",
			Updated: true,
			Links:   BuildDailyLogbookStatusLinks(baseURL, responseID, false),
		}

		log.Info(logger.LogDailyLogbookDeactivateOK, "logbook_id", logbookUUID)
		h.Response.SuccessWithData(c, domain.MsgDailyLogbookDeactivateOK, response)
	}
}
