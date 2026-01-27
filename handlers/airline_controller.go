package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAirlineByID godoc
// @Summary Get airline by ID
// @Description Retrieves airline information by its unique identifier
// @Tags airlines
// @Accept json
// @Produce json
// @Param id path string true "Airline ID (encoded or UUID)"
// @Security BearerAuth
// @Success 200 {object} middleware.APIResponse{data=AirlineResponse} "Airline retrieved successfully"
// @Failure 400 {object} middleware.APIResponse "Invalid ID format"
// @Failure 401 {object} middleware.APIResponse "Unauthorized"
// @Failure 404 {object} middleware.APIResponse "Airline not found"
// @Failure 500 {object} middleware.APIResponse "Internal server error"
// @Router /airlines/{id} [get]
func (h *handler) GetAirlineByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		if idParam == "" {
			Logger.Warn(logger.LogAirlineGet, "error", "empty id parameter")
			c.Error(domain.ErrInvalidID)
			return
		}

		// Try to decode the ID (hashids), if fails assume it's a raw UUID
		resolvedID := idParam
		encodedID := idParam
		if uuid, err := h.DecodeID(idParam); err == nil {
			resolvedID = uuid
		}

		airline, err := h.AirlineInteractor.GetAirlineByID(c.Request.Context(), resolvedID)
		if err != nil {
			switch err {
			case domain.ErrAirlineNotFound:
				c.Error(domain.ErrAirlineNotFound)
			default:
				c.Error(domain.ErrInternalServer)
			}
			return
		}

		// Encode the ID for response if not already encoded
		if resolvedID != idParam {
			encodedID = idParam
		} else {
			if enc, encErr := h.EncodeID(airline.ID); encErr == nil {
				encodedID = enc
			} else {
				encodedID = airline.ID
			}
		}

		response := FromDomainAirline(airline, encodedID)
		h.Response.SuccessWithData(c, domain.MsgAirlineGetOK, response)
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
		log := Logger.WithTraceID(traceID)

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
