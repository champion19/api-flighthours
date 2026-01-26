package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDHeader = "X-Request-ID"

	RequestIDKey = "request_id"

	TraceIDKey = "traceID"
)

type contextKey string
const traceIDContextKey contextKey = "traceID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(RequestIDKey, requestID)

		ctx:=context.WithValue(c.Request.Context(), traceIDContextKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

func GetTraceIDFromContext(ctx context.Context) string {

		if traceID, ok := ctx.Value(traceIDContextKey).(string); ok {
			return traceID
		}
		return ""
	}


