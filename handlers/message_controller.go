package handlers

import (
	domain "github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

// CreateMessage godoc
// @Summary      Create a new system message
// @Description  Creates a new message for the centralized messaging system
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        message body MessageRequest true "Message data"
// @Success      201  {object}  MessageCreatedResponse
// @Failure      400  {object}  middleware.ErrorResponse "Invalid request data"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      409  {object}  middleware.ErrorResponse "Message code already exists"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /messages [post]
func (h handler) CreateMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

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

		messageRequest.Sanitize()

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

		log.Success(logger.LogMessageCreateOK,
			"id", result.ID,
			"encoded_id", encodedID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMessageCreated, response)
	}
}

// UpdateMessage godoc
// @Summary      Update a system message
// @Description  Updates an existing message in the system
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Message ID (obfuscated)"
// @Param        message body MessageRequest true "Message data"
// @Success      200  {object}  MessageUpdatedResponse
// @Failure      400  {object}  middleware.ErrorResponse "Invalid request data"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Message not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /messages/{id} [put]
func (h handler) UpdateMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Info(logger.LogMessageUpdate,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		inputID := c.Param("id")
		uuid, responseID := h.resolveID(inputID)
		if uuid == "" {
			log.Error(logger.LogMessageInvalidID,
				"input_id", inputID,
				"client_ip", c.ClientIP())
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
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

		messageRequest.Sanitize()

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
		links := BuildMessageUpdatedLinks(baseURL, responseID)

		response := MessageUpdatedResponse{
			Links: links,
		}

		log.Success(logger.LogMessageUpdateOK,
			"id", result.ID,
			"code", result.Code,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMessageUpdated, response)
	}
}

// DeleteMessage godoc
// @Summary      Delete a system message
// @Description  Permanently deletes a message from the system
// @Tags         Messages
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Message ID (obfuscated)"
// @Success      200  {object}  middleware.APIResponse "Message deleted"
// @Failure      400  {object}  middleware.ErrorResponse "Invalid ID"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Message not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /messages/{id} [delete]
func (h handler) DeleteMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Info(logger.LogMessageDelete,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		inputID := c.Param("id")
		uuid, _ := h.resolveID(inputID)
		if uuid == "" {
			log.Error(logger.LogMessageInvalidID,
				"input_id", inputID,
				"client_ip", c.ClientIP())
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		log.Info(logger.LogMessageDeleteProcessing, "id", uuid)

		err := h.MessageInteractor.DeleteMessage(c, uuid)
		if err != nil {
			log.Error(logger.LogMessageDeleteError,
				"id", uuid,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		log.Success(logger.LogMessageDeleteOK,
			"id", uuid,
			"client_ip", c.ClientIP())

		h.Response.Success(c, domain.MsgMessageDeleted)
	}
}

// GetMessageByID godoc
// @Summary      Get a message by ID
// @Description  Returns a specific system message by its ID
// @Tags         Messages
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Message ID (obfuscated)"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  middleware.ErrorResponse "Invalid ID"
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      404  {object}  middleware.ErrorResponse "Message not found"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /messages/{id} [get]
func (h handler) GetMessageByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Debug(logger.LogMessageGet,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		inputID := c.Param("id")
		uuid, responseID := h.resolveID(inputID)
		if uuid == "" {
			log.Error(logger.LogMessageInvalidID,
				"input_id", inputID,
				"client_ip", c.ClientIP())
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
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

		baseURL := GetBaseURL(c)
		response := ToMessageResponse(message, responseID)
		response.Links = BuildMessageLinks(baseURL, responseID)

		log.Debug(logger.LogMessageGetOK,
			"id", uuid,
			"code", message.Code,
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

// ListMessages godoc
// @Summary      List all system messages
// @Description  Returns a list of all system messages with optional filters
// @Tags         Messages
// @Produce      json
// @Security     BearerAuth
// @Param        module query string false "Filter by module"
// @Param        type query string false "Filter by type (ERROR, WARNING, INFO, SUCCESS)"
// @Param        category query string false "Filter by category"
// @Param        active query string false "Filter by active status (true/false)"
// @Success      200  {object}  MessageListResponse
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /messages [get]
func (h handler) ListMessages() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Debug(logger.LogMessageList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		filters := buildMessageFilters(c)

		messages, err := h.MessageInteractor.ListMessages(c, filters)
		if err != nil {
			log.Error(logger.LogMessageListError,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToMessageListResponse(messages, h.EncodeID)

		for i := range response.Messages {
			response.Messages[i].Links = BuildMessageLinks(baseURL, response.Messages[i].ID)
		}
		response.Links = BuildMessageListLinks(baseURL)

		log.Debug(logger.LogMessageListOK,
			"count", len(messages),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}

func buildMessageFilters(c *gin.Context) map[string]interface{} {
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

	parseActiveFilter(c, filters)
	return filters
}

func parseActiveFilter(c *gin.Context, filters map[string]interface{}) {
	active := c.Query("active")
	switch active {
	case "true", "1":
		filters["active"] = true
	case "false", "0":
		filters["active"] = false
	}
}

// ReloadMessageCache godoc
// @Summary      Reload message cache
// @Description  Reloads the message cache from the database
// @Tags         Messages
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  CacheReloadResponse
// @Failure      401  {object}  middleware.ErrorResponse "Not authenticated"
// @Failure      500  {object}  middleware.ErrorResponse "Internal server error"
// @Router       /messages/cache/reload [post]
func (h handler) ReloadMessageCache() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := log.WithTraceID(traceID)

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
			Message:     "Caché de mensajes recargado exitosamente desde la base de datos",
		}

		log.Success(logger.LogMessageCacheReloadSuccess,
			"before_count", beforeCount,
			"after_count", afterCount,
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
