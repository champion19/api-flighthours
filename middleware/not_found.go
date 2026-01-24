package middleware

import (
	"net/http"

	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// NotFoundHandler logs 404 errors with trace ID and returns a JSON response
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID for trace correlation
		traceID := GetRequestID(c)
		log := log.WithTraceID(traceID)

		// Log the 404 with full context
		log.Warn(logger.LogMiddlewareNotFound,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent())


		// Return JSON response
		c.JSON(http.StatusNotFound, gin.H{

			"success": false,
			"code":    "404_NOT_FOUND",
			"message": "Endpoint no encontrado",
			"path":    c.Request.URL.Path,
		})


	}
}
