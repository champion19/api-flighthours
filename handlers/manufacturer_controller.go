package handlers

import(
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)




// GetManufacturerByID godoc
// @Summary      Get manufacturer by ID
// @Description  Returns manufacturer information by ID (accepts both UUID and obfuscated ID)
// @Tags         Manufacturers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Manufacturer ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /manufacturers/{id} [get]
func (h *handler) GetManufacturerByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogManufacturerGet, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		manufacturerUUID, responseID := h.resolveID(inputID)
		if manufacturerUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Get manufacturer from interactor
		manufacturer, err := h.ManufacturerInteractor.GetManufacturerByID(c.Request.Context(), manufacturerUUID)
		if err != nil {
			log.Error(logger.LogManufacturerGetError, "manufacturer_id", manufacturerUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrManufacturerNotFound {
				h.Response.Error(c, domain.MsgManufacturerNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainManufacturer(manufacturer, responseID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildManufacturerLinks(baseURL, responseID)

		log.Success(logger.LogManufacturerGetOK, manufacturer.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgManufacturerGetOK, response)
	}
}

// ListManufacturers godoc
// @Summary      List all manufacturers
// @Description  Returns a list of all manufacturers
// @Tags         Manufacturers
// @Produce      json
// @Success      200  {object}  ManufacturerListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /manufacturers [get]
func (h *handler) ListManufacturers() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogManufacturerList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		manufacturers, err := h.ManufacturerInteractor.ListManufacturers(c.Request.Context())
		if err != nil {
			log.Error(logger.LogManufacturerListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToManufacturerListResponse(manufacturers, h.EncodeID, baseURL)

		log.Debug(logger.LogManufacturerListOK,
			"count", len(manufacturers),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
