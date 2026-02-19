package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAirlineByID godoc
// @Summary      Get airline by ID
// @Description  Returns airline information by ID (accepts both UUID and obfuscated ID)
// @Tags         Airlines
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airlines/{id} [get]
func (h *handler) GetAirlineByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirlineGet, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		airlineUUID, responseID := h.resolveID(inputID)
		if airlineUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Get airline from interactor
		airline, err := h.AirlineInteractor.GetAirlineByID(c.Request.Context(), airlineUUID)
		if err != nil {
			log.Error(logger.LogAirlineGetError, "airline_id", airlineUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineNotFound {
				h.Response.Error(c, domain.MsgAirlineNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainAirline(airline, responseID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineLinks(baseURL, responseID)

		log.Success(logger.LogAirlineGetOK, airline.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirlineGetOK, response)
	}
}

// ActivateAirline godoc
// @Summary      Activate airline
// @Description  Sets airline status to active (accepts both UUID and obfuscated ID)
// @Tags         Airlines
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airlines/{id}/activate [patch]
func (h *handler) ActivateAirline() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirlineActivate, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		airlineUUID, responseID := h.resolveID(inputID)
		if airlineUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Activate airline via interactor
		if err := h.AirlineInteractor.ActivateAirline(c.Request.Context(), airlineUUID); err != nil {
			log.Error(logger.LogAirlineActivateError, "airline_id", airlineUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineNotFound {
				h.Response.Error(c, domain.MsgAirlineNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirlineStatusResponse{
			ID:      responseID,
			Status:  "active",
			Updated: true,
		}

		// Build HATEOAS links (isActive=true, muestra link para deactivate)
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineStatusLinks(baseURL, responseID, true)

		log.Success(logger.LogAirlineActivateOK, "airline_id", airlineUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineActivateOK, response)
	}
}

// DeactivateAirline godoc
// @Summary      Deactivate airline
// @Description  Sets airline status to inactive (accepts both UUID and obfuscated ID)
// @Tags         Airlines
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airlines/{id}/deactivate [patch]
func (h *handler) DeactivateAirline() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirlineDeactivate, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		airlineUUID, responseID := h.resolveID(inputID)
		if airlineUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Deactivate airline via interactor
		if err := h.AirlineInteractor.DeactivateAirline(c.Request.Context(), airlineUUID); err != nil {
			log.Error(logger.LogAirlineDeactivateError, "airline_id", airlineUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineNotFound {
				h.Response.Error(c, domain.MsgAirlineNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirlineStatusResponse{
			ID:      responseID,
			Status:  "inactive",
			Updated: true,
		}

		// Build HATEOAS links (isActive=false, muestra link para activate)
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineStatusLinks(baseURL, responseID, false)

		log.Success(logger.LogAirlineDeactivateOK, "airline_id", airlineUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineDeactivateOK, response)
	}
}

// ListAirlines godoc
// @Summary      List all airlines
// @Description  Returns a list of all airlines with optional status filter
// @Tags         Airlines
// @Produce      json
// @Param        status query string false "Filter by status (true for active, false for inactive)"
// @Success      200  {object}  AirlineListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /airlines [get]
func (h *handler) ListAirlines() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Debug(logger.LogAirlineList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Parse query parameters for filters
		filters := make(map[string]interface{})
		if status := c.Query("status"); status != "" {
			if status == "true" || status == "1" || status == "active" {
				filters["status"] = true
			} else if status == "false" || status == "0" || status == "inactive" {
				filters["status"] = false
			}
		}

		airlines, err := h.AirlineInteractor.ListAirlines(c.Request.Context(), filters)
		if err != nil {
			log.Error(logger.LogAirlineListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAirlineListResponse(airlines, h.EncodeID, baseURL)

		log.Debug(logger.LogAirlineListOK,
			"count", len(airlines),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
