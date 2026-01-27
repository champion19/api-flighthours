package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// CreateMessage godoc
// @Summary      Crear nuevo mensaje del sistema
// @Description  Crea un nuevo mensaje para el sistema de mensajes centralizados
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body handlers.MessageRequest true "Datos del mensaje a crear"
// @Success      201 {object} middleware.APIResponse{data=handlers.MessageCreatedResponse} "Mensaje creado exitosamente"
// @Failure      400 {object} middleware.APIResponse "Error de validación"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      409 {object} middleware.APIResponse "Código de mensaje duplicado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /messages [post]
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
		message.SetID()

		result, err := h.MessageInteractor.CreateMessage(c, message)
		if err != nil {
			log.Error(logger.LogMessageCreateError,
				"code", messageRequest.Code,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		encodedID, err := h.IDEncoder.Encode(result.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.ID, err)
			return
		}
		baseURL := GetBaseURL(c)
		links := BuildMessageCreatedLinks(baseURL, encodedID)
		SetLocationHeader(c, baseURL, "messages", encodedID)

		response := MessageCreatedResponse{
			ID:    encodedID,
			Links: links,
		}

		log.Success(logger.LogMessageCreatedSuccess,
			"id", result.ID,
			"encoded_id", encodedID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMessageCreated, response)
	}
}

// UpdateMessage godoc
// @Summary      Actualizar mensaje existente
// @Description  Actualiza un mensaje del sistema por su ID ofuscado
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID ofuscado del mensaje"
// @Param        request body handlers.MessageRequest true "Datos actualizados del mensaje"
// @Success      200 {object} middleware.APIResponse{data=handlers.MessageUpdatedResponse} "Mensaje actualizado exitosamente"
// @Failure      400 {object} middleware.APIResponse "Error de validación o ID inválido"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      404 {object} middleware.APIResponse "Mensaje no encontrado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /messages/{id} [put]
func (h handler) UpdateMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMessageUpdate,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

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

		baseURL := GetBaseURL(c)
		links := BuildMessageUpdatedLinks(baseURL, encodedID)

		response := MessageUpdatedResponse{
			Links: links,
		}

		log.Success(logger.LogMessageUpdatedSuccess,
			"id", result.ID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMessageUpdated, response)
	}
}

// DeleteMessage godoc
// @Summary      Eliminar mensaje
// @Description  Elimina un mensaje del sistema por su ID ofuscado
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID ofuscado del mensaje"
// @Success      200 {object} middleware.APIResponse "Mensaje eliminado exitosamente"
// @Failure      400 {object} middleware.APIResponse "ID inválido"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      404 {object} middleware.APIResponse "Mensaje no encontrado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /messages/{id} [delete]
func (h handler) DeleteMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMessageDelete,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

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

		log.Success(logger.LogMessageDeletedSuccess,
			"id", uuid,
			"client_ip", c.ClientIP())

		h.Response.Success(c, domain.MsgMessageDeleted)
	}
}

// GetMessageByID godoc
// @Summary      Obtener mensaje por ID
// @Description  Obtiene un mensaje específico por su ID ofuscado
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID ofuscado del mensaje"
// @Success      200 {object} handlers.MessageResponse "Mensaje encontrado"
// @Failure      400 {object} middleware.APIResponse "ID inválido"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      404 {object} middleware.APIResponse "Mensaje no encontrado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /messages/{id} [get]
func (h handler) GetMessageByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogMessageGet,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

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

		encodedIDForResponse, err := h.IDEncoder.Encode(message.ID)
		if err != nil {
			h.HandleIDEncodingError(c, message.ID, err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToMessageResponse(message)
		response.ID = encodedIDForResponse
		response.Links = BuildMessageLinks(baseURL, encodedIDForResponse)

		log.Debug(logger.LogMessageGetOK,
			"id", uuid,
			"code", message.Code,
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// ListMessages godoc
// @Summary      Listar todos los mensajes
// @Description  Lista todos los mensajes del sistema con filtros opcionales
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        module query string false "Filtrar por módulo"
// @Param        type query string false "Filtrar por tipo (SUCCESS, ERROR, WARNING, INFO)"
// @Param        category query string false "Filtrar por categoría"
// @Param        active query boolean false "Filtrar por estado activo"
// @Success      200 {object} handlers.MessageListResponse "Lista de mensajes"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      500 {object} middleware.APIResponse "Error interno del servidor"
// @Router       /messages [get]
func (h handler) ListMessages() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogMessageList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

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
// @Summary      Recargar cache de mensajes
// @Description  Recarga todos los mensajes desde la base de datos al cache en memoria
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} handlers.CacheReloadResponse "Cache recargado exitosamente"
// @Failure      401 {object} middleware.APIResponse "No autenticado"
// @Failure      500 {object} middleware.APIResponse "Error al recargar cache"
// @Router       /messages/cache/reload [post]
func (h handler) ReloadMessageCache() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMessageCacheReloading,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		beforeCount := h.MessagingCache.MessageCount()

		err := h.MessagingCache.ReloadMessages(c.Request.Context())
		if err != nil {
			log.Error(logger.LogMessageCacheReloadError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(domain.ErrInternalServer)
			return
		}

		afterCount := h.MessagingCache.MessageCount()

		response := CacheReloadResponse{
			Success:     true,
			BeforeCount: beforeCount,
			AfterCount:  afterCount,
			Message:     logger.LogMessageCacheReloadedMsg,
		}

		log.Success(logger.LogMessageCacheReloadSuccess,
			"before_count", beforeCount,
			"after_count", afterCount,
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
