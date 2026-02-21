package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

const errEmptyIDParam = "empty id parameter"

// GetAirportByID godoc
// @Summary      Get airport by ID
// @Description  Returns airport information by ID (accepts both UUID and obfuscated ID)
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airport ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airports/{id} [get]
func (h *handler) GetAirportByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", errEmptyIDParam, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportGet, "input_id", inputID, "client_ip", c.ClientIP())

		airportUUID, responseID := h.resolveID(inputID)
		if airportUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		airport, err := h.AirportInteractor.GetAirportByID(c.Request.Context(), airportUUID)
		if err != nil {
			log.Error(logger.LogAirportGetError, "airport_id", airportUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirportNotFound {
				h.Response.Error(c, domain.MsgAirportNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainAirport(airport, responseID)

		baseURL := GetBaseURL(c)
		response.Links = BuildAirportLinks(baseURL, responseID)

		log.Success(logger.LogAirportGetOK, airport.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirportGetOK, response)
	}
}

// ActivateAirport godoc
// @Summary      Activate airport
// @Description  Sets airport status to active (accepts both UUID and obfuscated ID)
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airport ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airports/{id}/activate [patch]
func (h *handler) ActivateAirport() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", errEmptyIDParam, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportActivate, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		airportUUID, responseID := h.resolveID(inputID)
		if airportUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Activate airport via interactor
		if err := h.AirportInteractor.ActivateAirport(c.Request.Context(), airportUUID); err != nil {
			log.Error(logger.LogAirportActivateError, "airport_id", airportUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirportNotFound {
				h.Response.Error(c, domain.MsgAirportNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirportStatusResponse{
			ID:      responseID,
			Status:  "active",
			Updated: true,
		}

		// Build HATEOAS links (isActive=true, muestra link para deactivate)
		baseURL := GetBaseURL(c)
		response.Links = BuildAirportStatusLinks(baseURL, responseID, true)

		log.Success(logger.LogAirportActivateOK, "airport_id", airportUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirportActivateOK, response)
	}
}

// DeactivateAirport godoc
// @Summary      Deactivate airport
// @Description  Sets airport status to inactive (accepts both UUID and obfuscated ID)
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airport ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airports/{id}/deactivate [patch]
func (h *handler) DeactivateAirport() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", errEmptyIDParam, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportDeactivate, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		airportUUID, responseID := h.resolveID(inputID)
		if airportUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		if err := h.AirportInteractor.DeactivateAirport(c.Request.Context(), airportUUID); err != nil {
			log.Error(logger.LogAirportDeactivateError, "airport_id", airportUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirportNotFound {
				h.Response.Error(c, domain.MsgAirportNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirportStatusResponse{
			ID:      responseID,
			Status:  "inactive",
			Updated: true,
		}

		baseURL := GetBaseURL(c)
		response.Links = BuildAirportStatusLinks(baseURL, responseID, false)

		log.Success(logger.LogAirportDeactivateOK, "airport_id", airportUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirportDeactivateOK, response)
	}
}

// ListAirports godoc
// @Summary      List all airports
// @Description  Returns a list of all airports with optional status filter
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Param        status  query     string  false  "Filter by status (true/false, active/inactive)"
// @Success      200  {object}  AirportListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /airports [get]
func (h *handler) ListAirports() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Debug(logger.LogAirportList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		filters := make(map[string]interface{})
		if status := c.Query("status"); status != "" {
			if status == "true" || status == "1" || status == "active" {
				filters["status"] = true
			} else if status == "false" || status == "0" || status == "inactive" {
				filters["status"] = false
			}
		}

		airports, err := h.AirportInteractor.ListAirports(c.Request.Context(), filters)
		if err != nil {
			log.Error(logger.LogAirportListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToAirportListResponse(airports, h.EncodeID, baseURL)

		log.Debug(logger.LogAirportListOK,
			"count", len(airports),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// GetAirportsByType godoc
// @Summary      Get airports by type
// @Description  Returns all airports of a specific type. No new tables needed - queries airport.airport_type field.
// @Tags         Airport Types
// @Accept       json
// @Produce      json
// @Param        airport_type   path      string  true  "Airport type (e.g., INTERNACIONAL, NACIONAL)"
// @Success      200  {object}  AirportListResponse
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airport-types/{airport_type} [get]
func (h *handler) GetAirportsByType() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		airportType := c.Param("airport_type")
		if airportType == "" {
			log.Error(logger.LogAirportTypeGetError, "error", "empty airport_type parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportTypeGet, "airport_type", airportType, "client_ip", c.ClientIP())

		airports, err := h.AirportInteractor.GetAirportsByType(c.Request.Context(), airportType)
		if err != nil {
			log.Error(logger.LogAirportTypeGetError, "airport_type", airportType, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAirportTypeNotFound)
			return
		}
		baseURL := GetBaseURL(c)
		response := ToAirportListResponse(airports, h.EncodeID, baseURL)

		response.Links = append(response.Links, Link{
			Rel:  "airports",
			Href: baseURL + "/airports",
		})

		log.Success(logger.LogAirportTypeGetOK, "airport_type", airportType, "count", len(airports), "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirportTypeGetOK, response)
	}
}
