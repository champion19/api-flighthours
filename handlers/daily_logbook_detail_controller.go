package handlers

import(
 "github.com/gin-gonic/gin"
 "github.com/champion19/api-flighthours/middleware"
 "github.com/champion19/api-flighthours/platform/logger"
 "github.com/champion19/api-flighthours/core/interactor/services/domain"
)

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
func (h *handler) GetDailyLogbookDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailGet, "input_id", inputID)

		// Resolve ID (supports both UUID and obfuscated ID)
		detailUUID, responseID := h.resolveID(inputID)
		if detailUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailGetError, "error", "invalid ID")
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
		encodedAircraftID, _ := h.EncodeID(detail.LicensePlateID)

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
func (h *handler) CreateDailyLogbookDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailCreate, "logbook_id", inputID)

		// Get authenticated user
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		// Resolve logbook ID
		logbookUUID, _ := h.resolveID(inputID)
		if logbookUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid logbook ID")
			h.Response.Error(c, domain.MsgFlightInvalidLogbook)
			return
		}

		// Verify ownership
		if err := h.DailyLogbookDetailInteractor.VerifyLogbookOwnership(c.Request.Context(), logbookUUID, employee.ID); err != nil {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "unauthorized")
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgFlightUnauthorized)
				return
			}
			h.Response.Error(c, domain.MsgFlightInvalidLogbook)
			return
		}

		// Parse and sanitize request
		var req CreateDailyLogbookDetailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}
		req.Sanitize()

		// Resolve airline_route_id
		routeUUID, _ := h.resolveID(req.AirlineRouteID)
		if routeUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid route ID")
			h.Response.Error(c, domain.MsgFlightInvalidRoute)
			return
		}
		req.AirlineRouteID = routeUUID

		// Resolve aircraft_registration_id
		aircraftUUID, _ := h.resolveID(req.LicensePlateID)
		if aircraftUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid aircraft ID")
			h.Response.Error(c, domain.MsgFlightInvalidAircraft)
			return
		}
		req.LicensePlateID = aircraftUUID

		// Validate pilot role
		if !domain.IsValidPilotRole(req.PilotRole) {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid pilot role")
			h.Response.Error(c, domain.MsgValFieldFormat)
			return
		}

		// Validate approach type if provided
		if req.ApproachType != nil && !domain.IsValidApproachType(*req.ApproachType) {
			log.Warn(logger.LogDailyLogbookDetailCreateError, "error", "invalid approach type")
			h.Response.Error(c, domain.MsgValFieldFormat)
			return
		}

		// Convert to domain
		detail := ToDomainDailyLogbookDetail(logbookUUID, req)
		detail.SetID()

		// Set employee logbook ID
		detail.EmployeeLogbookID = &employee.ID

		// Create detail
		if err := h.DailyLogbookDetailInteractor.CreateDailyLogbookDetail(c.Request.Context(), traceID, detail); err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", err)
			if err == domain.ErrFlightInvalidLogbook {
				h.Response.Error(c, domain.MsgFlightInvalidLogbook)
				return
			}
			if err == domain.ErrFlightInvalidRoute {
				h.Response.Error(c, domain.MsgFlightInvalidRoute)
				return
			}
			if err == domain.ErrFlightInvalidAircraft {
				h.Response.Error(c, domain.MsgFlightInvalidAircraft)
				return
			}
			if err == domain.ErrFlightInvalidTimeSequence {
				h.Response.Error(c, domain.MsgFlightInvalidTimeSequence)
				return
			}
			h.Response.Error(c, domain.MsgFlightSaveError)
			return
		}

		// Refetch to get denormalized data
		createdDetail, err := h.DailyLogbookDetailInteractor.GetDailyLogbookDetailByID(c.Request.Context(), traceID, detail.ID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailCreateError, "error", err)
			h.Response.Error(c, domain.MsgFlightSaveError)
			return
		}

		// Encode IDs for response
		encodedID, _ := h.EncodeID(detail.ID)
		encodedLogbookID, _ := h.EncodeID(logbookUUID)
		encodedRouteID, _ := h.EncodeID(req.AirlineRouteID)
		encodedAircraftID, _ := h.EncodeID(req.LicensePlateID)

		// Build response
		response := FromDomainDailyLogbookDetail(createdDetail, encodedID, encodedLogbookID, encodedRouteID, encodedAircraftID)
		response.Links = BuildDailyLogbookDetailLinks(c, encodedID)

		log.Info(logger.LogDailyLogbookDetailCreateOK, "id", detail.ID)
		h.Response.SuccessWithData(c, domain.MsgFlightCreated, response)
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
func (h *handler) UpdateDailyLogbookDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)
		inputID := c.Param("id")

		log.Info(logger.LogDailyLogbookDetailUpdate, "input_id", inputID)

		// Get authenticated user
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		// Resolve detail ID
		detailUUID, responseID := h.resolveID(inputID)
		if detailUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "invalid ID")
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Verify ownership via detail's logbook
		ownerID, err := h.DailyLogbookDetailInteractor.GetDetailLogbookOwner(c.Request.Context(), detailUUID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
			if err == domain.ErrFlightNotFound {
				h.Response.Error(c, domain.MsgFlightNotFound)
				return
			}
			h.Response.Error(c, domain.MsgFlightUpdateError)
			return
		}
		if ownerID != employee.ID {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "unauthorized")
			h.Response.Error(c, domain.MsgFlightUnauthorized)
			return
		}

		// Parse and sanitize request
		var req UpdateDailyLogbookDetailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}
		req.Sanitize()

		// Resolve airline_route_id
		routeUUID, _ := h.resolveID(req.AirlineRouteID)
		if routeUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "invalid route ID")
			h.Response.Error(c, domain.MsgFlightInvalidRoute)
			return
		}
		req.AirlineRouteID = routeUUID

		// Resolve aircraft_registration_id
		aircraftUUID, _ := h.resolveID(req.LicensePlateID)
		if aircraftUUID == "" {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "invalid aircraft ID")
			h.Response.Error(c, domain.MsgFlightInvalidAircraft)
			return
		}
		req.LicensePlateID = aircraftUUID

		// Validate pilot role
		if !domain.IsValidPilotRole(req.PilotRole) {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "invalid pilot role")
			h.Response.Error(c, domain.MsgValFieldFormat)
			return
		}

		// Validate approach type if provided
		if req.ApproachType != nil && !domain.IsValidApproachType(*req.ApproachType) {
			log.Warn(logger.LogDailyLogbookDetailUpdateError, "error", "invalid approach type")
			h.Response.Error(c, domain.MsgValFieldFormat)
			return
		}

		// Convert to domain
		detail := ToDomainDailyLogbookDetailUpdate(detailUUID, req)

		// Update detail
		if err := h.DailyLogbookDetailInteractor.UpdateDailyLogbookDetail(c.Request.Context(), traceID, detail); err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
			if err == domain.ErrFlightNotFound {
				h.Response.Error(c, domain.MsgFlightNotFound)
				return
			}
			if err == domain.ErrFlightInvalidTimeSequence {
				h.Response.Error(c, domain.MsgFlightInvalidTimeSequence)
				return
			}
			h.Response.Error(c, domain.MsgFlightUpdateError)
			return
		}

		// Refetch to get denormalized data
		updatedDetail, err := h.DailyLogbookDetailInteractor.GetDailyLogbookDetailByID(c.Request.Context(), traceID, detailUUID)
		if err != nil {
			log.Error(logger.LogDailyLogbookDetailUpdateError, "error", err)
			h.Response.Error(c, domain.MsgFlightUpdateError)
			return
		}

		// Encode IDs for response
		encodedLogbookID, _ := h.EncodeID(updatedDetail.DailyLogbookID)
		encodedRouteID, _ := h.EncodeID(updatedDetail.AirlineRouteID)
		encodedAircraftID, _ := h.EncodeID(updatedDetail.LicensePlateID)

		// Build response
		response := FromDomainDailyLogbookDetail(updatedDetail, responseID, encodedLogbookID, encodedRouteID, encodedAircraftID)
		response.Links = BuildDailyLogbookDetailLinks(c, responseID)

		log.Info(logger.LogDailyLogbookDetailUpdateOK, "id", detailUUID)
		h.Response.SuccessWithData(c, domain.MsgFlightUpdated, response)
	}
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
func (h *handler) ListDailyLogbookDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)
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

		// Verify ownership
		if err := h.DailyLogbookDetailInteractor.VerifyLogbookOwnership(c.Request.Context(), logbookUUID, employee.ID); err != nil {
			log.Warn(logger.LogDailyLogbookDetailListError, "error", "unauthorized")
			if err == domain.ErrFlightUnauthorized {
				h.Response.Error(c, domain.MsgFlightUnauthorized)
				return
			}
			h.Response.Error(c, domain.MsgFlightInvalidLogbook)
			return
		}

		// Get details
		details, err := h.DailyLogbookDetailInteractor.ListDailyLogbookDetailsByLogbook(c.Request.Context(), traceID, logbookUUID)
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
			encodedAircraftID, _ := h.EncodeID(d.LicensePlateID)

			response := FromDomainDailyLogbookDetail(&d, encodedID, encodedLogbookID, encodedRouteID, encodedAircraftID)
			response.Links = BuildDailyLogbookDetailLinks(c, encodedID)
			responses = append(responses, response)
		}

		log.Info(logger.LogDailyLogbookDetailListOK, "count", len(responses))
		h.Response.SuccessWithData(c, domain.MsgFlightListOK, responses)
	}
}
