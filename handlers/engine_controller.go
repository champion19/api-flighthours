package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetEngineByID godoc
// @Summary      Get engine by ID
// @Description  Returns engine information by ID (accepts both UUID and obfuscated ID)
// @Tags         Engines
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Engine ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /engines/{id} [get]
func (h *handler) GetEngineByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogEngineGet, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		engineUUID, responseID := h.resolveID(inputID)
		if engineUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Get engine from interactor
		engine, err := h.EngineInteractor.GetEngineByID(c.Request.Context(), engineUUID)
		if err != nil {
			log.Error(logger.LogEngineGetError, "engine_id", engineUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrEngineNotFound {
				h.Response.Error(c, domain.MsgEngineNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainEngine(engine, responseID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildEngineLinks(baseURL, responseID)

		log.Success(logger.LogEngineGetOK, engine.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgEngineGetOK, response)
	}
}

// ListEngines godoc
// @Summary      List all engines
// @Description  Returns a list of all engine types
// @Tags         Engines
// @Produce      json
// @Success      200  {object}  EngineListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /engines [get]
func (h *handler) ListEngines() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogEngineList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		engines, err := h.EngineInteractor.ListEngines(c.Request.Context())
		if err != nil {
			log.Error(logger.LogEngineListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToEngineListResponse(engines, h.EncodeID, baseURL)

		log.Debug(logger.LogEngineListOK,
			"count", len(engines),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
