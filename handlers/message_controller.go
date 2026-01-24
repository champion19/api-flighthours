package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
	"github.com/champion19/flighthours-api/middleware"
)


// CreateMessage godoc
// @Summary Crear un nuevo mensaje del sistema
// @Description Crea un nuevo mensaje del sistema con tipo, categoría y contenido. Los mensajes se utilizan para mostrar información, advertencias y errores a los usuarios.
// @Tags messages
// @Accept json
// @Produce json
// @Param message body MessageRequest true "Datos del mensaje a crear"
// @Success 201 {object} MessageCreatedResponse "Mensaje creado exitosamente"
// @Failure 400 {object} middleware.ErrorResponse "Datos de entrada inválidos"
// @Failure 409 {object} middleware.ErrorResponse "El código del mensaje ya existe"
// @Failure 500 {object} middleware.ErrorResponse "Error interno del servidor"
// @Router /messages [post]
func (h handler) CreateMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMessageCreate,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var messageRequest MessageRequest
		if err := c.ShouldBindJSON(&messageRequest); err != nil {
			log.Error(logger.LogMiddlewareJSONParseError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		log.Info(logger.LogMessageCreateProcessing,
			"code", messageRequest.Code,
			"type", messageRequest.Type)

		message := messageRequest.ToDomain()
		message.SetID() // Generate UUID

		result, err := h.MessageInteractor.CreateMessage(c, message)
		if err != nil {
			log.Error(logger.LogMessageCreateError,
				"code", messageRequest.Code,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Encode UUID for public API
		encodedID, err := h.IDEncoder.Encode(result.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.ID, err)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildMessageCreatedLinks(baseURL, encodedID)

		// Set Location header
		SetLocationHeader(c, baseURL, "messages", encodedID)

		response := MessageCreatedResponse{
			ID:    encodedID,
			Links: links,
		}

		log.Success("Mensaje creado exitosamente",
			"id", result.ID,
			"encoded_id", encodedID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMessageCreated, response)
	}
}

// UpdateMessage godoc
// @Summary Actualizar un mensaje del sistema
// @Description Actualiza un mensaje del sistema existente. Permite modificar tipo, categoría, contenido y estado activo.
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "ID del mensaje (encoded)"
// @Param message body MessageRequest true "Datos del mensaje a actualizar"
// @Success 200 {object} MessageUpdatedResponse "Mensaje actualizado exitosamente"
// @Failure 400 {object} middleware.ErrorResponse "Datos de entrada inválidos"
// @Failure 404 {object} middleware.ErrorResponse "Mensaje no encontrado"
// @Failure 500 {object} middleware.ErrorResponse "Error interno del servidor"
// @Router /messages/{id} [put]
func (h handler) UpdateMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMessageUpdate,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Get encoded ID from URL parameter and decode to UUID
		encodedID := c.Param("id")
		uuid, err := h.IDEncoder.Decode(encodedID)
		if err != nil {
			log.Error(logger.LogMessageInvalidID,
				"encoded_id", encodedID,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidID)
			return
		}

		var messageRequest MessageRequest
		if err := c.ShouldBindJSON(&messageRequest); err != nil {
			log.Error(logger.LogMiddlewareJSONParseError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		log.Info(logger.LogMessageUpdateProcessing,
			"id", uuid,
			"code", messageRequest.Code)

		message := messageRequest.ToDomain()
		message.ID = uuid

		result, err := h.MessageInteractor.UpdateMessage(c, message)
		if err != nil {
			log.Error(logger.LogMessageUpdateError,
				"id", uuid,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildMessageUpdatedLinks(baseURL, encodedID)

		response := MessageUpdatedResponse{
			Links: links,
		}

		log.Success("Mensaje actualizado exitosamente",
			"id", result.ID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMessageUpdated, response)
	}
}

// DeleteMessage godoc
// @Summary Eliminar un mensaje del sistema
// @Description Elimina un mensaje del sistema de forma permanente. Esta acción no se puede deshacer.
// @Tags messages
// @Produce json
// @Param id path string true "ID del mensaje (encoded)"
// @Success 200 {object} MessageDeletedResponse "Mensaje eliminado exitosamente"
// @Failure 400 {object} middleware.ErrorResponse "ID inválido"
// @Failure 404 {object} middleware.ErrorResponse "Mensaje no encontrado"
// @Failure 500 {object} middleware.ErrorResponse "Error interno del servidor"
// @Router /messages/{id} [delete]
func (h handler) DeleteMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMessageDelete,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Get encoded ID from URL parameter and decode to UUID
		encodedID := c.Param("id")
		uuid, err := h.IDEncoder.Decode(encodedID)
		if err != nil {
			log.Error(logger.LogMessageInvalidID,
				"encoded_id", encodedID,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidID)
			return
		}

		log.Info(logger.LogMessageDeleteProcessing, "id", uuid)

		err = h.MessageInteractor.DeleteMessage(c, uuid)
		if err != nil {
			log.Error(logger.LogMessageDeleteError,
				"id", uuid,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		log.Success("Mensaje eliminado exitosamente",
			"id", uuid,
			"client_ip", c.ClientIP())

		h.Response.Success(c, domain.MsgMessageDeleted)
	}
}

// GetMessageByID godoc
// @Summary Obtener un mensaje por ID
// @Description Obtiene los detalles de un mensaje del sistema específico por su ID.
// @Tags messages
// @Produce json
// @Param id path string true "ID del mensaje (encoded)"
// @Success 200 {object} MessageResponse "Mensaje encontrado exitosamente"
// @Failure 400 {object} middleware.ErrorResponse "ID inválido"
// @Failure 404 {object} middleware.ErrorResponse "Mensaje no encontrado"
// @Failure 500 {object} middleware.ErrorResponse "Error interno del servidor"
// @Router /messages/{id} [get]
func (h handler) GetMessageByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogMessageGet,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Get encoded ID from URL parameter and decode to UUID
		encodedID := c.Param("id")
		uuid, err := h.IDEncoder.Decode(encodedID)
		if err != nil {
			log.Error(logger.LogMessageInvalidID,
				"encoded_id", encodedID,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidID)
			return
		}

		message, err := h.MessageInteractor.GetMessageByID(c, uuid)
		if err != nil {
			log.Error(logger.LogMessageGetError,
				"id", uuid,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Encode UUID for response
		encodedIDForResponse, err := h.IDEncoder.Encode(message.ID)
		if err != nil {
			h.HandleIDEncodingError(c, message.ID, err)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToMessageResponse(message)
		response.ID = encodedIDForResponse // Use encoded ID in response
		response.Links = BuildMessageLinks(baseURL, encodedIDForResponse)

		log.Debug(logger.LogMessageGetOK,
			"id", uuid,
			"code", message.Code,
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// ListMessages godoc
// @Summary Listar mensajes del sistema
// @Description Obtiene una lista de mensajes del sistema con filtros opcionales. Permite filtrar por módulo, tipo, categoría y estado activo.
// @Tags messages
// @Produce json
// @Param module query string false "Filtrar por módulo (ej: users, flights, bookings)"
// @Param type query string false "Filtrar por tipo (ERROR, WARNING, INFO, SUCCESS)"
// @Param category query string false "Filtrar por categoría (usuario_final, sistema, validacion)"
// @Param active query string false "Filtrar por estado activo (true, false)"
// @Success 200 {object} MessageListResponse "Lista de mensajes obtenida exitosamente"
// @Failure 500 {object} middleware.ErrorResponse "Error interno del servidor"
// @Router /messages [get]
func (h handler) ListMessages() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogMessageList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Parse query parameters for filters
		filters := make(map[string]interface{})
		if module := c.Query("module"); module != "" {
			filters["module"] = module
		}
		if msgType := c.Query("type"); msgType != "" {
			filters["type"] = msgType
		}
		if category := c.Query("category"); category != "" {
			filters["category"] = category
		}
		if active := c.Query("active"); active != "" {
			if active == "true" || active == "1" {
				filters["active"] = true
			} else if active == "false" || active == "0" {
				filters["active"] = false
			}
		}

		messages, err := h.MessageInteractor.ListMessages(c, filters)
		if err != nil {
			log.Error(logger.LogMessageListError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Encode UUIDs for each message in response
		baseURL := GetBaseURL(c)
		response := ToMessageListResponse(messages)
		for i := range response.Messages {
			encodedID, err := h.IDEncoder.Encode(messages[i].ID)
			if err != nil {
				h.HandleIDEncodingError(c, messages[i].ID, err)
				return
			}
			response.Messages[i].ID = encodedID
			response.Messages[i].Links = BuildMessageLinks(baseURL, encodedID)
		}
		response.Links = BuildMessageListLinks(baseURL)

		log.Debug(logger.LogMessageListOK,
			"count", len(messages),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// ReloadMessageCache godoc
// @Summary Recargar caché de mensajes
// @Description Recarga el caché de mensajes desde la base de datos. Útil después de hacer cambios masivos a mensajes.
// @Tags messages
// @Produce json
// @Success 200 {object} CacheReloadResponse "Caché recargado exitosamente"
// @Failure 500 {object} middleware.ErrorResponse "Error interno del servidor"
// @Router /messages/reload [post]
func (h handler) ReloadMessageCache() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info("Recargando caché de mensajes",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Obtener el conteo antes del reload
		beforeCount := h.MessagingCache.MessageCount()

		// Recargar el caché desde BD
		err := h.MessagingCache.ReloadMessages(c.Request.Context())
		if err != nil {
			log.Error("Error al recargar caché de mensajes",
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInternalServer)
			return
		}

		// Obtener el conteo después del reload
		afterCount := h.MessagingCache.MessageCount()

		response := CacheReloadResponse{
			Success:     true,
			BeforeCount: beforeCount,
			AfterCount:  afterCount,
			Message:     "Caché de mensajes recargado exitosamente desde la base de datos",
		}

		log.Success("Caché de mensajes recargado exitosamente",
			"before_count", beforeCount,
			"after_count", afterCount,
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
