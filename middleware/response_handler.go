package middleware

import (
	"net/http"
   messagingCache "github.com/champion19/flighthours-api/platform/cache/messaging"
	"github.com/gin-gonic/gin"
)



type ResponseHandler struct {
	cache *messagingCache.MessageCache
}

func NewResponseHandler(cache *messagingCache.MessageCache) *ResponseHandler {
	return &ResponseHandler{
		cache: cache,
	}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}


func (h *ResponseHandler) Error(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Code:    code,
			Message: "Unknown error",
		})
		return
	}

	c.JSON(status, APIResponse{
		Success: false,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

func (h *ResponseHandler) ErrorWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Code:    code,
			Message: "Unknown error",
			Data:    data,
		})
		return
	}

	c.JSON(status, APIResponse{
		Success: false,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

func (h *ResponseHandler) Success(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "Operation successful",
		})
		return
	}

	c.JSON(status, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

func (h *ResponseHandler) SuccessWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)
	status := h.cache.GetHTTPStatus(code)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "Operation successful",
			Data:    data,
		})
		return
	}

	c.JSON(status, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

func (h *ResponseHandler) Warning(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System warning",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

func (h *ResponseHandler) WarningWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System warning",
			Data:    data,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

func (h *ResponseHandler) Info(c *gin.Context, code string, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System information",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
	})
}

func (h *ResponseHandler) InfoWithData(c *gin.Context, code string, data interface{}, params ...string) {
	msg := h.cache.GetMessageResponse(code, params...)

	if msg == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Code:    code,
			Message: "System information",
			Data:    data,
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    msg.Code,
		Message: msg.Content,
		Data:    data,
	})
}

func (h *ResponseHandler) DataOnly(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}
