package handlers

import (
	"context"
	"net/http"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetTailNumberByPlate godoc
// @Summary      Get aircraft registration by license plate number
// @Description  Returns aircraft registration information by its plate number (e.g. HK-5432)
// @Tags         License Plates
// @Accept       json
// @Produce      json
// @Param        plate   path      string  true  "License Plate number (e.g. HK-5432)"
// @Success      200  {object}  TailNumberResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tail-numbers/{plate} [get]
func (h *handler) GetTailNumberByPlate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := log.WithTraceID(traceID)

		plate := c.Param("plate")
		if plate == "" {
			_ = c.Error(domain.ErrInvalidID)
			return
		}

		registration, err := h.TailNumberInteractor.GetTailNumberByPlate(ctx, plate)
		if err != nil {
			log.Error(logger.LogTailNumberGetError, "tail_number", plate, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		encodedID, err := h.EncodeID(registration.ID)
		if err != nil {
			h.HandleIDEncodingError(c, registration.ID, err)
			return
		}

		encodedModelID, _ := h.EncodeID(registration.AircraftModelID)
		encodedAirlineID, _ := h.EncodeID(registration.AirlineID)

		baseURL := GetBaseURL(c)
		response := FromDomainTailNumber(registration, encodedID, encodedModelID, encodedAirlineID)
		response.Links = BuildTailNumberLinks(baseURL, encodedID)

		log.Success(logger.LogTailNumberGetOK, registration.ToLogger())
		c.JSON(http.StatusOK, response)
	}
}

// ListTailNumbers godoc
// @Summary      List all aircraft registrations
// @Description  Returns a list of all aircraft registrations with optional filters
// @Tags         License Plates
// @Produce      json
// @Param        airline_id    query string false "Filter by airline ID"
// @Param        tail_number query string false "Filter by license plate"
// @Success      200  {object}  TailNumberListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /tail-numbers [get]
func (h *handler) ListTailNumbers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := log.WithTraceID(traceID)

		filters := make(map[string]interface{})

		if tailNumber := c.Query("tail_number"); tailNumber != "" {
			filters["tail_number"] = tailNumber
		} else if airlineID := c.Query("airline_id"); airlineID != "" {
			resolvedAirlineID, _ := h.resolveID(airlineID)
			if resolvedAirlineID != "" {
				filters["airline_id"] = resolvedAirlineID
			}
		}

		registrations, err := h.TailNumberInteractor.ListTailNumbers(ctx, filters)
		if err != nil {
			log.Error(logger.LogTailNumberListError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToTailNumberListResponse(registrations, h.EncodeID, baseURL)

		log.Success(logger.LogTailNumberListOK, "count", len(registrations))
		c.JSON(http.StatusOK, response)
	}
}

// CreateTailNumber godoc
// @Summary      Create a new aircraft registration
// @Description  Creates a new aircraft registration (license plate must be unique)
// @Tags         License Plates
// @Accept       json
// @Produce      json
// @Param        request body CreateTailNumberRequest true "License Plate data"
// @Success      201  {object}  TailNumberResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tail-numbers [post]
func (h *handler) CreateTailNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := log.WithTraceID(traceID)

		var req CreateTailNumberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogRegJSONBindError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		resolvedModelID, _ := h.resolveID(req.AircraftModelID)
		if resolvedModelID == "" {
			log.Warn(logger.LogTailNumberInvalidModelID, "id", req.AircraftModelID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrTailNumberInvalidModel)
			return
		}
		req.AircraftModelID = resolvedModelID

		resolvedAirlineID, _ := h.resolveID(req.AirlineID)
		if resolvedAirlineID == "" {
			log.Warn(logger.LogTailNumberInvalidAirlineID, "id", req.AirlineID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrTailNumberInvalidAirline)
			return
		}
		req.AirlineID = resolvedAirlineID

		registration := req.ToDomain()

		if err := h.TailNumberInteractor.CreateTailNumber(ctx, registration); err != nil {
			log.Error(logger.LogTailNumberCreateError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		createdRegistration, err := h.TailNumberInteractor.GetTailNumberByID(ctx, registration.ID)
		if err != nil {
			log.Warn(logger.LogTailNumberFetchAfterWriteError, "id", registration.ID, "error", err)

			createdRegistration = &registration
		}

		encodedID, err := h.EncodeID(registration.ID)
		if err != nil {
			h.HandleIDEncodingError(c, registration.ID, err)
			return
		}

		encodedModelID, _ := h.EncodeID(createdRegistration.AircraftModelID)
		encodedAirlineID, _ := h.EncodeID(createdRegistration.AirlineID)

		baseURL := GetBaseURL(c)
		response := FromDomainTailNumber(createdRegistration, encodedID, encodedModelID, encodedAirlineID)
		response.Links = BuildTailNumberCreatedLinks(baseURL, encodedID)

		SetLocationHeader(c, baseURL, "tail-numbers", encodedID)
		log.Success(logger.LogTailNumberCreateOK, createdRegistration.ToLogger())
		c.JSON(http.StatusCreated, response)
	}
}

// UpdateTailNumber godoc
// @Summary      Update an existing aircraft registration
// @Description  Updates an aircraft registration by ID (accepts both UUID and obfuscated ID)
// @Tags         License Plates
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "License Plate ID (obfuscated ID)"
// @Param        request body UpdateTailNumberRequest true "License Plate data"
// @Success      200  {object}  TailNumberResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /tail-numbers/{id} [put]
func (h *handler) UpdateTailNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			_ = c.Error(domain.ErrInvalidID)
			return
		}

		registrationUUID, responseID := h.resolveID(inputID)
		if registrationUUID == "" {
			h.checkDuplicateOnInvalidID(c, ctx, inputID, log)
			return
		}

		var req UpdateTailNumberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogRegJSONBindError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		resolvedModelID, _ := h.resolveID(req.AircraftModelID)
		if resolvedModelID == "" {
			log.Warn(logger.LogTailNumberInvalidModelID, "id", req.AircraftModelID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrTailNumberInvalidModel)
			return
		}
		req.AircraftModelID = resolvedModelID

		resolvedAirlineID, _ := h.resolveID(req.AirlineID)
		if resolvedAirlineID == "" {
			log.Warn(logger.LogTailNumberInvalidAirlineID, "id", req.AirlineID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrTailNumberInvalidAirline)
			return
		}
		req.AirlineID = resolvedAirlineID

		if h.isDuplicateData(ctx, registrationUUID, req) {
			log.Warn(logger.LogTailNumberUpdateError, "id", registrationUUID, "reason", "duplicate_data", "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrTailNumberDuplicatePlate)
			return
		}

		registration := req.ToDomain(registrationUUID)

		if err := h.TailNumberInteractor.UpdateTailNumber(ctx, registration); err != nil {
			log.Error(logger.LogTailNumberUpdateError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		updatedRegistration, err := h.TailNumberInteractor.GetTailNumberByID(ctx, registration.ID)
		if err != nil {
			log.Warn(logger.LogTailNumberFetchAfterWriteError, "id", registration.ID, "error", err)

			updatedRegistration = &registration
		}

		encodedModelID, _ := h.EncodeID(updatedRegistration.AircraftModelID)
		encodedAirlineID, _ := h.EncodeID(updatedRegistration.AirlineID)

		baseURL := GetBaseURL(c)
		response := FromDomainTailNumber(updatedRegistration, responseID, encodedModelID, encodedAirlineID)
		response.Links = BuildTailNumberLinks(baseURL, responseID)

		log.Success(logger.LogTailNumberUpdateOK, updatedRegistration.ToLogger())
		c.JSON(http.StatusOK, response)
	}
}

func (h *handler) checkDuplicateOnInvalidID(c *gin.Context, ctx context.Context, inputID string, log logger.Logger) {
	var checkReq UpdateTailNumberRequest
	if c.ShouldBindJSON(&checkReq) == nil {
		checkReq.Sanitize()
		if checkReq.TailNumber != "" {
			filters := map[string]interface{}{"tail_number": checkReq.TailNumber}
			existing, _ := h.TailNumberInteractor.ListTailNumbers(ctx, filters)
			if len(existing) > 0 {
				log.Warn(logger.LogTailNumberUpdateError, "id", inputID, "tail_number", checkReq.TailNumber, "reason", "duplicate", "client_ip", c.ClientIP())
				_ = c.Error(domain.ErrTailNumberDuplicatePlate)
				return
			}
		}
	}
	log.Warn(logger.LogTailNumberNotFound, "id", inputID, "client_ip", c.ClientIP())
	_ = c.Error(domain.ErrTailNumberNotFound)
}

func (h *handler) isDuplicateData(ctx context.Context, registrationUUID string, req UpdateTailNumberRequest) bool {
	currentReg, err := h.TailNumberInteractor.GetTailNumberByID(ctx, registrationUUID)
	if err != nil || currentReg == nil {
		return false
	}
	return currentReg.TailNumber == req.TailNumber &&
		currentReg.AircraftModelID == req.AircraftModelID &&
		currentReg.AirlineID == req.AirlineID
}
