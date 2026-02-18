package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAircraftModelByID godoc
// @Summary      Get aircraft model by ID
// @Description  Returns aircraft model information by ID (accepts both UUID and obfuscated ID)
// @Tags         Aircraft Models
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Aircraft Model ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-models/{id} [get]
func (h *handler) GetAircraftModelByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAircraftModelGet, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		modelUUID, responseID := h.resolveID(inputID)
		if modelUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Get aircraft model from interactor
		model, err := h.AircraftModelInteractor.GetAircraftModelByID(c.Request.Context(), modelUUID)
		if err != nil {
			log.Error(logger.LogAircraftModelGetError, "aircraft_model_id", modelUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAircraftModelNotFound {
				h.Response.Error(c, domain.MsgAircraftModelNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainAircraftModel(model, responseID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAircraftModelLinks(baseURL, responseID)

		log.Success(logger.LogAircraftModelGetOK, model.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAircraftModelGetOK, response)
	}
}

// ListAircraftModels godoc
// @Summary      List all aircraft models
// @Description  Returns a list of all aircraft models with optional engine type filter
// @Tags         Aircraft Models
// @Produce      json
// @Param        engine_type query string false "Filter by engine type (e.g., JET, TUR)"
// @Success      200  {object}  AircraftModelListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-models [get]
func (h *handler) ListAircraftModels() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Debug(logger.LogAircraftModelList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Parse query parameters for filters
		filters := make(map[string]interface{})
		if engineType := c.Query("engine_type"); engineType != "" {
			filters["engine_type"] = engineType
		}

		models, err := h.AircraftModelInteractor.ListAircraftModels(c.Request.Context(), filters)
		if err != nil {
			log.Error(logger.LogAircraftModelListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAircraftModelListResponse(models, h.EncodeID, baseURL)

		log.Debug(logger.LogAircraftModelListOK,
			"count", len(models),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// GetAircraftModelsByFamily godoc
// @Summary      Get aircraft family by family name
// @Description  Returns aircraft family information (family and manufacturer) for a specific family (HU32)
// @Tags         Aircraft Families
// @Produce      json
// @Param        family   path      string  true  "Aircraft Family (e.g., A320, B737)"
// @Success      200  {object}  AircraftFamilyListResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-families/{family} [get]
func (h *handler) GetAircraftModelsByFamily() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		familyID := c.Param("family")
		if familyID == "" {
			log.Error(logger.LogAircraftFamilyGetError, "error", "empty family parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAircraftFamilyGet, "family", familyID, "client_ip", c.ClientIP())

		// Get aircraft models by family from interactor
		models, err := h.AircraftModelInteractor.GetAircraftModelsByFamily(c.Request.Context(), familyID)
		if err != nil {
			log.Error(logger.LogAircraftFamilyGetError, "family", familyID, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAircraftFamilyGetErr)
			return
		}

		// If no models found for this family, return 404
		if len(models) == 0 {
			log.Warn(logger.LogAircraftFamilyNotFound, "family", familyID, "message", "no models found", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAircraftFamilyNotFound)
			return
		}

		// Convert to family response (only family + manufacturer)
		baseURL := GetBaseURL(c)
		response := ToAircraftFamilyListResponse(models, baseURL)

		log.Success(logger.LogAircraftFamilyGetOK, "family", familyID, "count", response.Total, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAircraftFamilyGetOK, response)
	}
}

// ActivateAircraftModel godoc
// @Summary      Activate aircraft model
// @Description  Sets aircraft model status to active (HU42 - accepts obfuscated ID only)
// @Tags         Aircraft Models
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Aircraft Model ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-models/{id}/activate [patch]
func (h *handler) ActivateAircraftModel() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAircraftModelActivate, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		modelUUID, responseID := h.resolveID(inputID)
		if modelUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Activate aircraft model via interactor
		if err := h.AircraftModelInteractor.ActivateAircraftModel(c.Request.Context(), modelUUID); err != nil {
			log.Error(logger.LogAircraftModelActivateError, "aircraft_model_id", modelUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAircraftModelNotFound {
				h.Response.Error(c, domain.MsgAircraftModelNotFound)
				return
			}
			h.Response.Error(c, domain.MsgAircraftModelActivateErr)
			return
		}

		response := AircraftModelStatusResponse{
			ID:      responseID,
			Status:  "active",
			Updated: true,
		}

		// Build HATEOAS links (isActive=true, shows link for deactivate)
		baseURL := GetBaseURL(c)
		response.Links = BuildAircraftModelStatusLinks(baseURL, responseID, true)

		log.Success(logger.LogAircraftModelActivateOK, "aircraft_model_id", modelUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAircraftModelActivateOK, response)
	}
}

// DeactivateAircraftModel godoc
// @Summary      Deactivate aircraft model
// @Description  Sets aircraft model status to inactive (HU41 - accepts obfuscated ID only)
// @Tags         Aircraft Models
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Aircraft Model ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-models/{id}/deactivate [patch]
func (h *handler) DeactivateAircraftModel() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAircraftModelDeactivate, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		modelUUID, responseID := h.resolveID(inputID)
		if modelUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Deactivate aircraft model via interactor
		if err := h.AircraftModelInteractor.DeactivateAircraftModel(c.Request.Context(), modelUUID); err != nil {
			log.Error(logger.LogAircraftModelDeactivateError, "aircraft_model_id", modelUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAircraftModelNotFound {
				h.Response.Error(c, domain.MsgAircraftModelNotFound)
				return
			}
			h.Response.Error(c, domain.MsgAircraftModelDeactivateErr)
			return
		}

		response := AircraftModelStatusResponse{
			ID:      responseID,
			Status:  "inactive",
			Updated: true,
		}

		// Build HATEOAS links (isActive=false, shows link for activate)
		baseURL := GetBaseURL(c)
		response.Links = BuildAircraftModelStatusLinks(baseURL, responseID, false)

		log.Success(logger.LogAircraftModelDeactivateOK, "aircraft_model_id", modelUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAircraftModelDeactivateOK, response)
	}
}
