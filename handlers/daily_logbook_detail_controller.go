package handlers

import (
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

const errInvalidID = "invalid ID"

// GetDailyLogbookDetail retrieves a daily logbook detail by ID
// @Summary Get daily logbook detail by ID
// @Description Retrieves a specific daily logbook detail (flight segment)
// @Tags DailyLogbookDetails
// @Accept json
// @Produce json
// @Param id path string true "Detail ID (obfuscated or UUID)"
// @Success 200 {object} DailyLogbookDetailResponse
// @Failure 404 {object} middleware.APIResponse
// @Router /daily-logbook-details/{id} [get]
// @Security BearerAuth
func (h *handler) GetDailyLogbookDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailGet, "input_id", inputID)

		// Resolve ID (supports both UUID and obfuscated ID)
		detailUUID, responseID := h.resolveID(inputID)
		if detailUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailGetError, "error", errInvalidID)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Get detail
		detail, err := h.DailyLogbookDetailInteractor.GetDailyLogbookDetailByID(c.Request.Context(), traceID, detailUUID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailGetError, "error", err)
			if err == domain.ErrFlightNotFound {
				h.Response.Error(c, domain.MsgFlightNotFound)
				return
			}
			h.Response.Error(c, domain.MsgFlightGetErr)
			return
		}

		// Encode related IDs
		encodedLogbookID, _ := h.EncodeID(detail.DailyLogbookID)
		encodedRouteID, _ := h.EncodeID(detail.AirlineRouteID)
		encodedAircraftID, _ := h.EncodeID(detail.TailNumberID)

		// Build response
		response := FromDomainDailyLogbookDetail(detail, responseID, encodedLogbookID, encodedRouteID, encodedAircraftID)
		response.Links = BuildDailyLogbookDetailLinks(c, responseID)

		log.Info(logger.LogDailyLogbookDetailGetOK, "id", detailUUID)
		h.Response.SuccessWithData(c, domain.MsgFlightGetOK, response)
	}
}

// CreateDailyLogbookDetail creates a new detail under a logbook
// @Summary Create daily logbook detail
// @Description Creates a new flight segment under a daily logbook
// @Tags DailyLogbookDetails
// @Accept json
// @Produce json
// @Param id path string true "Logbook ID (obfuscated or UUID)"
// @Param body body CreateDailyLogbookDetailRequest true "Detail data"
// @Success 201 {object} DailyLogbookDetailResponse
// @Failure 400 {object} middleware.APIResponse
// @Failure 403 {object} middleware.APIResponse
// @Router /daily-logbooks/{id}/details [post]
// @Security BearerAuth
func (h *handler) CreateDailyLogbookDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailCreate, "logbook_id", inputID)

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		logbookUUID, _ := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid logbook ID")
			h.Response.Error(c, domain.MsgFlightInvalidLogbook)
			return
		}

		var req CreateDailyLogbookDetailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}
		req.Sanitize()

		routeUUID, _ := h.resolveID(req.AirlineRouteID)
		if routeUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid route ID")
			h.Response.Error(c, domain.MsgFlightInvalidRoute)
			return
		}
		req.AirlineRouteID = routeUUID

		var tailNumberUUID string
		if req.TailNumberID == "" {
			// No tail number given for this flight — fall back to the one already
			// set on the parent daily logbook (captured once in "New Logbook Entry").
			logbook, lerr := h.DailyLogbookInteractor.GetDailyLogbookByID(c.Request.Context(), logbookUUID, employee.ID)
			if lerr == nil && logbook != nil && logbook.TailNumberID != nil && *logbook.TailNumberID != "" {
				tailNumberUUID = *logbook.TailNumberID
			}
		} else {
			tailNumberUUID, _ = h.resolveID(req.TailNumberID)
		}
		if tailNumberUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid tail number ID")
			h.Response.Error(c, domain.MsgFlightInvalidTailNumber)
			return
		}
		req.TailNumberID = tailNumberUUID

		detail := ToDomainDailyLogbookDetail(logbookUUID, req)
		detail.SetID()
		detail.EmployeeLogbookID = &employee.ID

		if err := domain.ValidateApproachFields(detail.ApproachCategory, detail.ApproachSubtype, detail.Autoland); err != nil {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", err)
			h.Response.Error(c, domain.MsgFlightInvalidApproach)
			return
		}

		if err := h.DailyLogbookDetailInteractor.CreateDailyLogbookDetail(c.Request.Context(), traceID, detail, employee.ID); err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", err)
			h.Response.Error(c, mapCreateError(err))
			return
		}

		createdDetail, err := h.DailyLogbookDetailInteractor.GetDailyLogbookDetailByID(c.Request.Context(), traceID, detail.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", err)
			h.Response.Error(c, domain.MsgFlightSaveError)
			return
		}

		encodedID, _ := h.EncodeID(detail.ID)
		encodedLogbookID, _ := h.EncodeID(logbookUUID)
		encodedRouteID, _ := h.EncodeID(req.AirlineRouteID)
		encodedAircraftID, _ := h.EncodeID(req.TailNumberID)

		response := FromDomainDailyLogbookDetail(createdDetail, encodedID, encodedLogbookID, encodedRouteID, encodedAircraftID)
		response.Links = BuildDailyLogbookDetailLinks(c, encodedID)

		log.Info(logger.LogDailyLogbookDetailCreateOK, "id", detail.ID)
		h.Response.SuccessWithData(c, domain.MsgFlightCreated, response)
	}
}

func mapCreateError(err error) string {
	switch err {
	case domain.ErrFlightUnauthorized:
		return domain.MsgFlightUnauthorized
	case domain.ErrFlightInvalidLogbook:
		return domain.MsgFlightInvalidLogbook
	case domain.ErrFlightInvalidRoute:
		return domain.MsgFlightInvalidRoute
	case domain.ErrFlightInvalidTailNumber:
		return domain.MsgFlightInvalidTailNumber
	case domain.ErrFlightInvalidTimeSequence:
		return domain.MsgFlightInvalidTimeSequence
	case domain.ErrFlightDuplicate:
		return domain.MsgFlightDuplicate
	default:
		return domain.MsgFlightSaveError
	}
}

// UpdateDailyLogbookDetail updates an existing detail
// @Summary Update daily logbook detail
// @Description Updates a flight segment
// @Tags DailyLogbookDetails
// @Accept json
// @Produce json
// @Param id path string true "Detail ID (obfuscated or UUID)"
// @Param body body UpdateDailyLogbookDetailRequest true "Detail data"
// @Success 200 {object} DailyLogbookDetailResponse
// @Failure 400 {object} middleware.APIResponse
// @Failure 403 {object} middleware.APIResponse
// @Failure 404 {object} middleware.APIResponse
// @Router /daily-logbook-details/{id} [put]
// @Security BearerAuth
func (h *handler) UpdateDailyLogbookDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailUpdate, "input_id", inputID)

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		detailUUID, responseID := h.resolveID(inputID)
		if detailUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", errInvalidID)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		var req UpdateDailyLogbookDetailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}
		req.Sanitize()

		routeUUID, _ := h.resolveID(req.AirlineRouteID)
		if routeUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "invalid route ID")
			h.Response.Error(c, domain.MsgFlightInvalidRoute)
			return
		}
		req.AirlineRouteID = routeUUID

		aircraftUUID, _ := h.resolveID(req.TailNumberID)
		if aircraftUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "invalid tail number ID")
			h.Response.Error(c, domain.MsgFlightInvalidTailNumber)
			return
		}
		req.TailNumberID = aircraftUUID

		detail := ToDomainDailyLogbookDetailUpdate(detailUUID, req)

		if err := domain.ValidateApproachFields(detail.ApproachCategory, detail.ApproachSubtype, detail.Autoland); err != nil {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", err)
			h.Response.Error(c, domain.MsgFlightInvalidApproach)
			return
		}

		if err := h.DailyLogbookDetailInteractor.UpdateDailyLogbookDetail(c.Request.Context(), traceID, detail, employee.ID); err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
			h.Response.Error(c, mapUpdateError(err))
			return
		}

		updatedDetail, err := h.DailyLogbookDetailInteractor.GetDailyLogbookDetailByID(c.Request.Context(), traceID, detailUUID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
			h.Response.Error(c, domain.MsgFlightUpdateError)
			return
		}

		encodedLogbookID, _ := h.EncodeID(updatedDetail.DailyLogbookID)
		encodedRouteID, _ := h.EncodeID(updatedDetail.AirlineRouteID)
		encodedAircraftID, _ := h.EncodeID(updatedDetail.TailNumberID)

		response := FromDomainDailyLogbookDetail(updatedDetail, responseID, encodedLogbookID, encodedRouteID, encodedAircraftID)
		response.Links = BuildDailyLogbookDetailLinks(c, responseID)

		log.Info(logger.LogDailyLogbookDetailUpdateOK, "id", detailUUID)
		h.Response.SuccessWithData(c, domain.MsgFlightUpdated, response)
	}
}

func mapUpdateError(err error) string {
	switch err {
	case domain.ErrFlightUnauthorized:
		return domain.MsgFlightUnauthorized
	case domain.ErrFlightNotFound:
		return domain.MsgFlightNotFound
	case domain.ErrFlightInvalidTimeSequence:
		return domain.MsgFlightInvalidTimeSequence
	default:
		return domain.MsgFlightUpdateError
	}
}

// DeleteDailyLogbookDetail deletes a detail
// @Summary Delete daily logbook detail
// @Description Deletes a flight segment
// @Tags DailyLogbookDetails
// @Accept json
// @Produce json
// @Param id path string true "Detail ID (obfuscated or UUID)"
// @Success 200 {object} middleware.APIResponse
// @Failure 403 {object} middleware.APIResponse
// @Failure 404 {object} middleware.APIResponse
// @Router /daily-logbook-details/{id} [delete]
// @Security BearerAuth
func (h *handler) DeleteDailyLogbookDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailDelete, "input_id", inputID)

		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogDailyLogbookDetailDeleteError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		detailUUID, _ := h.resolveID(inputID)
		if detailUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailDeleteError, "error", errInvalidID)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		ownerID, err := h.DailyLogbookDetailInteractor.GetDetailLogbookOwner(c.Request.Context(), detailUUID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailDeleteError, "error", err)
			h.Response.Error(c, mapDeleteError(err))
			return
		}
		if ownerID != employee.ID {
			log.Warn(logger.LogDailyLogbookDetailDeleteError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		if err := h.DailyLogbookDetailInteractor.DeleteDailyLogbookDetail(c.Request.Context(), traceID, detailUUID); err != nil {
			log.Error(logger.LogDailyLogbookDetailDeleteError, "error", err)
			h.Response.Error(c, mapDeleteError(err))
			return
		}

		log.Info(logger.LogDailyLogbookDetailDeleteOK, "id", detailUUID)
		h.Response.Success(c, domain.MsgFlightDeleted)
	}
}

func mapDeleteError(err error) string {
	if err == domain.ErrFlightNotFound {
		return domain.MsgFlightNotFound
	}
	return domain.MsgFlightDeleteError
}

// ListDailyLogbookDetails lists all details for a logbook
// @Summary List daily logbook details
// @Description Lists all flight segments for a specific daily logbook
// @Tags DailyLogbookDetails
// @Accept json
// @Produce json
// @Param id path string true "Logbook ID (obfuscated or UUID)"
// @Success 200 {array} DailyLogbookDetailResponse
// @Failure 404 {object} middleware.APIResponse
// @Router /daily-logbooks/{id}/details [get]
// @Security BearerAuth
func (h *handler) ListDailyLogbookDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailList, "logbook_id", inputID)

		// Get authenticated user
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogDailyLogbookDetailListError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		// Resolve logbook ID
		logbookUUID, _ := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailListError, "error", "invalid logbook ID")
			h.Response.Error(c, domain.MsgFlightInvalidLogbook)
			return
		}

		// Get details (ownership is verified by the interactor)
		details, err := h.DailyLogbookDetailInteractor.ListDailyLogbookDetailsByLogbook(c.Request.Context(), traceID, logbookUUID, employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailListError, "error", err)
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgFlightUnauthorized)
				return
			}
			if err == domain.ErrFlightInvalidLogbook {
				h.Response.Error(c, domain.MsgFlightInvalidLogbook)
				return
			}
			h.Response.Error(c, domain.MsgFlightListError)
			return
		}

		// Build response
		var responses []DailyLogbookDetailResponse
		for _, d := range details {
			encodedID, _ := h.EncodeID(d.ID)
			encodedLogbookID, _ := h.EncodeID(d.DailyLogbookID)
			encodedRouteID, _ := h.EncodeID(d.AirlineRouteID)
			encodedAircraftID, _ := h.EncodeID(d.TailNumberID)

			response := FromDomainDailyLogbookDetail(&d, encodedID, encodedLogbookID, encodedRouteID, encodedAircraftID)
			response.Links = BuildDailyLogbookDetailLinks(c, encodedID)
			responses = append(responses, response)
		}

		log.Info(logger.LogDailyLogbookDetailListOK, "count", len(responses))
		h.Response.SuccessWithData(c, domain.MsgFlightListOK, responses)
	}
}

// ListMyFlights lists all flights for the authenticated employee
// @Summary List my flights
// @Description Lists all flight segments across all daily logbooks for the authenticated employee
// @Tags DailyLogbookDetails
// @Accept json
// @Produce json
// @Success 200 {array} DailyLogbookDetailResponse
// @Failure 401 {object} middleware.APIResponse
// @Router /employees/flights [get]
// @Security BearerAuth
func (h *handler) ListMyFlights() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Info(logger.LogDailyLogbookDetailList, "action", "list_my_flights")

		// Get authenticated user
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogDailyLogbookDetailListError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		// Get all flight details for this employee
		details, err := h.DailyLogbookDetailInteractor.ListDailyLogbookDetailsByEmployee(c.Request.Context(), traceID, employee.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailListError, "error", err)
			h.Response.Error(c, domain.MsgFlightListError)
			return
		}

		// Build response
		var responses []DailyLogbookDetailResponse
		for _, d := range details {
			encodedID, _ := h.EncodeID(d.ID)
			encodedLogbookID, _ := h.EncodeID(d.DailyLogbookID)
			encodedRouteID, _ := h.EncodeID(d.AirlineRouteID)
			encodedAircraftID, _ := h.EncodeID(d.TailNumberID)

			response := FromDomainDailyLogbookDetail(&d, encodedID, encodedLogbookID, encodedRouteID, encodedAircraftID)
			response.Links = BuildDailyLogbookDetailLinks(c, encodedID)
			responses = append(responses, response)
		}

		log.Info(logger.LogDailyLogbookDetailListOK, "employee_id", employee.ID, "count", len(responses))
		h.Response.SuccessWithData(c, domain.MsgFlightListOK, responses)
	}
}
