package handlers

import (
	"net/http"

	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetLicensePlateByID godoc
// @Summary      Get aircraft registration by ID
// @Description  Returns aircraft registration information by ID (accepts both UUID and obfuscated ID)
// @Tags         License Plates
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "License Plate ID (obfuscated ID)"
// @Success      200  {object}  LicensePlateResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /license-plates/{id} [get]
func (h *handler) GetLicensePlateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			_ = c.Error(domain.ErrInvalidID)
			return
		}

		registrationUUID, responseID := h.resolveID(inputID)
		if registrationUUID == "" {
			log.Warn(logger.LogLicensePlateNotFound, "id", inputID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrLicensePlateNotFound)
			return
		}

		registration, err := h.LicensePlateInteractor.GetLicensePlateByID(ctx, registrationUUID)
		if err != nil {
			log.Error(logger.LogLicensePlateGetError, "id", inputID, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		encodedModelID, _ := h.EncodeID(registration.AircraftModelID)
		encodedAirlineID, _ := h.EncodeID(registration.AirlineID)

		baseURL := GetBaseURL(c)
		response := FromDomainLicensePlate(registration, responseID, encodedModelID, encodedAirlineID)
		response.Links = BuildLicensePlateLinks(baseURL, responseID)

		log.Success(logger.LogLicensePlateGetOK, registration.ToLogger())
		c.JSON(http.StatusOK, response)
	}
}

// ListLicensePlates godoc
// @Summary      List all aircraft registrations
// @Description  Returns a list of all aircraft registrations with optional filters
// @Tags         License Plates
// @Produce      json
// @Param        airline_id    query string false "Filter by airline ID"
// @Param        license_plate query string false "Filter by license plate"
// @Success      200  {object}  LicensePlateListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /license-plates [get]
func (h *handler) ListLicensePlates() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		filters := make(map[string]interface{})

		if licensePlate := c.Query("license_plate"); licensePlate != "" {
			filters["license_plate"] = licensePlate
		} else if airlineID := c.Query("airline_id"); airlineID != "" {
			resolvedAirlineID, _ := h.resolveID(airlineID)
			if resolvedAirlineID != "" {
				filters["airline_id"] = resolvedAirlineID
			}
		}

		registrations, err := h.LicensePlateInteractor.ListLicensePlates(ctx, filters)
		if err != nil {
			log.Error(logger.LogLicensePlateListError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToLicensePlateListResponse(registrations, h.EncodeID, baseURL)

		log.Success(logger.LogLicensePlateListOK, "count", len(registrations))
		c.JSON(http.StatusOK, response)
	}
}

// CreateLicensePlate godoc
// @Summary      Create a new aircraft registration
// @Description  Creates a new aircraft registration (license plate must be unique)
// @Tags         License Plates
// @Accept       json
// @Produce      json
// @Param        request body CreateLicensePlateRequest true "License Plate data"
// @Success      201  {object}  LicensePlateResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /license-plates [post]
func (h *handler) CreateLicensePlate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		var req CreateLicensePlateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogRegJSONBindError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		resolvedModelID, _ := h.resolveID(req.AircraftModelID)
		if resolvedModelID == "" {
			log.Warn(logger.LogLicensePlateInvalidModelID, "id", req.AircraftModelID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrLicensePlateInvalidModel)
			return
		}
		req.AircraftModelID = resolvedModelID

		resolvedAirlineID, _ := h.resolveID(req.AirlineID)
		if resolvedAirlineID == "" {
			log.Warn(logger.LogLicensePlateInvalidAirlineID, "id", req.AirlineID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrLicensePlateInvalidAirline)
			return
		}
		req.AirlineID = resolvedAirlineID

		registration := req.ToDomain()

		if err := h.LicensePlateInteractor.CreateLicensePlate(ctx, registration); err != nil {
			log.Error(logger.LogLicensePlateCreateError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		createdRegistration, err := h.LicensePlateInteractor.GetLicensePlateByID(ctx, registration.ID)
		if err != nil {
			log.Warn(logger.LogLicensePlateFetchAfterWriteError, "id", registration.ID, "error", err)

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
		response := FromDomainLicensePlate(createdRegistration, encodedID, encodedModelID, encodedAirlineID)
		response.Links = BuildLicensePlateCreatedLinks(baseURL, encodedID)

		SetLocationHeader(c, baseURL, "license-plates", encodedID)
		log.Success(logger.LogLicensePlateCreateOK, createdRegistration.ToLogger())
		c.JSON(http.StatusCreated, response)
	}
}

// UpdateLicensePlate godoc
// @Summary      Update an existing aircraft registration
// @Description  Updates an aircraft registration by ID (accepts both UUID and obfuscated ID)
// @Tags         License Plates
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "License Plate ID (obfuscated ID)"
// @Param        request body UpdateLicensePlateRequest true "License Plate data"
// @Success      200  {object}  LicensePlateResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /license-plates/{id} [put]
func (h *handler) UpdateLicensePlate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			_ = c.Error(domain.ErrInvalidID)
			return
		}

		registrationUUID, responseID := h.resolveID(inputID)
		if registrationUUID == "" {
			// Before returning 404, check if the license plate from the body already exists
			var checkReq UpdateLicensePlateRequest
			if err := c.ShouldBindJSON(&checkReq); err == nil {
				checkReq.Sanitize()
				if checkReq.LicensePlate != "" {
					filters := map[string]interface{}{"license_plate": checkReq.LicensePlate}
					existing, _ := h.LicensePlateInteractor.ListLicensePlates(ctx, filters)
					if len(existing) > 0 {
						log.Warn(logger.LogLicensePlateUpdateError, "id", inputID, "license_plate", checkReq.LicensePlate, "reason", "duplicate", "client_ip", c.ClientIP())
						_ = c.Error(domain.ErrLicensePlateDuplicatePlate)
						return
					}
				}
			}
			log.Warn(logger.LogLicensePlateNotFound, "id", inputID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrLicensePlateNotFound)
			return
		}

		var req UpdateLicensePlateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogRegJSONBindError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		resolvedModelID, _ := h.resolveID(req.AircraftModelID)
		if resolvedModelID == "" {
			log.Warn(logger.LogLicensePlateInvalidModelID, "id", req.AircraftModelID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrLicensePlateInvalidModel)
			return
		}
		req.AircraftModelID = resolvedModelID

		resolvedAirlineID, _ := h.resolveID(req.AirlineID)
		if resolvedAirlineID == "" {
			log.Warn(logger.LogLicensePlateInvalidAirlineID, "id", req.AirlineID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrLicensePlateInvalidAirline)
			return
		}
		req.AirlineID = resolvedAirlineID

		// Check if the data is identical to what's already stored → duplicate
		currentReg, err := h.LicensePlateInteractor.GetLicensePlateByID(ctx, registrationUUID)
		if err == nil && currentReg != nil {
			if currentReg.LicensePlate == req.LicensePlate &&
				currentReg.AircraftModelID == req.AircraftModelID &&
				currentReg.AirlineID == req.AirlineID {
				log.Warn(logger.LogLicensePlateUpdateError, "id", registrationUUID, "reason", "duplicate_data", "client_ip", c.ClientIP())
				_ = c.Error(domain.ErrLicensePlateDuplicatePlate)
				return
			}
		}

		registration := req.ToDomain(registrationUUID)

		if err := h.LicensePlateInteractor.UpdateLicensePlate(ctx, registration); err != nil {
			log.Error(logger.LogLicensePlateUpdateError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		updatedRegistration, err := h.LicensePlateInteractor.GetLicensePlateByID(ctx, registration.ID)
		if err != nil {
			log.Warn(logger.LogLicensePlateFetchAfterWriteError, "id", registration.ID, "error", err)

			updatedRegistration = &registration
		}

		encodedModelID, _ := h.EncodeID(updatedRegistration.AircraftModelID)
		encodedAirlineID, _ := h.EncodeID(updatedRegistration.AirlineID)

		baseURL := GetBaseURL(c)
		response := FromDomainLicensePlate(updatedRegistration, responseID, encodedModelID, encodedAirlineID)
		response.Links = BuildLicensePlateLinks(baseURL, responseID)

		log.Success(logger.LogLicensePlateUpdateOK, updatedRegistration.ToLogger())
		c.JSON(http.StatusOK, response)
	}
}
