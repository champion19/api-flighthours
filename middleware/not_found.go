package middleware

import (
	"net/http"

	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
)

func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := GetRequestID(c)
		log := log.WithTraceID(traceID)

		log.Warn(logger.LogMiddlewareNotFound,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent())


		c.JSON(http.StatusNotFound, gin.H{

			"success": false,
			"code":    "404_NOT_FOUND",
			"message": "Endpoint no encontrado",
			"path":    c.Request.URL.Path,
		})


	}
}
