package handlers

import(
  domain	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"

)

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
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportGet, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		airportUUID, responseID := h.resolveID(inputID)
		if airportUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Get airport from interactor
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

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAirportLinks(baseURL, responseID)

		log.Success(logger.LogAirportGetOK, airport.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirportGetOK, response)
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
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
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

		// Deactivate airport via interactor
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

		// Build HATEOAS links (isActive=false, muestra link para activate)
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
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogAirportList,
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

		airports, err := h.AirportInteractor.ListAirports(c.Request.Context(), filters)
		if err != nil {
			log.Error(logger.LogAirportListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAirportListResponse(airports, h.EncodeID, baseURL)

		log.Debug(logger.LogAirportListOK,
			"count", len(airports),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// GetAirportsByCity godoc
// @Summary      Get airports by city (HU13 - Virtual Entity pattern)
// @Description  Returns all airports located in a specific city. No new tables needed - queries airport.city field.
// @Tags         Cities
// @Accept       json
// @Produce      json
// @Param        city_name   path      string  true  "City name (e.g., Bogota, Medellin)"
// @Success      200  {object}  AirportListResponse
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /cities/{city_name} [get]
func (h *handler) GetAirportsByCity() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		cityName := c.Param("city_name")
		if cityName == "" {
			log.Error(logger.LogAirportListError, "error", "empty city_name parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportList, "city_name", cityName, "client_ip", c.ClientIP())

		airports, err := h.AirportInteractor.GetAirportsByCity(c.Request.Context(), cityName)
		if err != nil {
			log.Error(logger.LogAirportListError, "city_name", cityName, "error", err, "client_ip", c.ClientIP())
			// If no rows found, city doesn't exist (no airports in that city)
			h.Response.Error(c, domain.MsgCityNotFound)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAirportListResponse(airports, h.EncodeID, baseURL)

		// Add city-specific links
		response.Links = append(response.Links, Link{
			Rel:  "airports",
			Href: baseURL + "/airports",
		})

		log.Success(logger.LogAirportListOK, "city_name", cityName, "count", len(airports), "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgCityGetOK, response)
	}
}

// GetAirportsByCountry godoc
// @Summary      Get airports by country (HU38 - Virtual Entity pattern)
// @Description  Returns all airports located in a specific country. No new tables needed - queries airport.country field.
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Param        country_name   path      string  true  "Country name (e.g., Colombia, Peru)"
// @Success      200  {object}  AirportListResponse
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /countries/{country_name} [get]
func (h *handler) GetAirportsByCountry() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		countryName := c.Param("country_name")
		if countryName == "" {
			log.Error(logger.LogAirportListError, "error", "empty country_name parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportList, "country_name", countryName, "client_ip", c.ClientIP())

		airports, err := h.AirportInteractor.GetAirportsByCountry(c.Request.Context(), countryName)
		if err != nil {
			log.Error(logger.LogAirportListError, "country_name", countryName, "error", err, "client_ip", c.ClientIP())
			// If no rows found, country doesn't exist (no airports in that country)
			h.Response.Error(c, domain.MsgCountryNotFound)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAirportListResponse(airports, h.EncodeID, baseURL)

		// Add country-specific links
		response.Links = append(response.Links, Link{
			Rel:  "airports",
			Href: baseURL + "/airports",
		})

		log.Success(logger.LogAirportListOK, "country_name", countryName, "count", len(airports), "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgCountryGetOK, response)
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
		log := Logger.WithTraceID(traceID)

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
			// If no rows found, airport type doesn't exist (no airports of that type)
			h.Response.Error(c, domain.MsgAirportTypeNotFound)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAirportListResponse(airports, h.EncodeID, baseURL)

		// Add airport type-specific links
		response.Links = append(response.Links, Link{
			Rel:  "airports",
			Href: baseURL + "/airports",
		})

		log.Success(logger.LogAirportTypeGetOK, "airport_type", airportType, "count", len(airports), "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirportTypeGetOK, response)
	}
}

