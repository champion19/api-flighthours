package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("uses existing request id header", func(t *testing.T) {
		r := gin.New()
		r.Use(RequestID())
		r.GET("/", func(c *gin.Context) {
			c.Status(200)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(RequestIDHeader, "abc")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Header().Get(RequestIDHeader) != "abc" {
			t.Fatalf("expected response header %s=abc, got %q", RequestIDHeader, w.Header().Get(RequestIDHeader))
		}
	})

	t.Run("generates request id when missing", func(t *testing.T) {
		r := gin.New()
		r.Use(RequestID())
		r.GET("/", func(c *gin.Context) {
			c.Status(200)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		id := w.Header().Get(RequestIDHeader)
		if id == "" {
			t.Fatalf("expected non-empty %s header", RequestIDHeader)
		}
	})
}

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns request id from context", func(t *testing.T) {
		r := gin.New()
		r.Use(RequestID())
		var capturedID string
		r.GET("/", func(c *gin.Context) {
			capturedID = GetRequestID(c)
			c.Status(200)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(RequestIDHeader, "test-request-123")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if capturedID != "test-request-123" {
			t.Errorf("expected 'test-request-123', got %q", capturedID)
		}
	})

	t.Run("returns empty string when not set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		id := GetRequestID(c)
		if id != "" {
			t.Errorf("expected empty string, got %q", id)
		}
	})
}

func TestGetTraceIDFromContext(t *testing.T) {
	t.Run("returns trace id when present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), traceIDContextKey, "trace-abc-123")
		id := GetTraceIDFromContext(ctx)
		if id != "trace-abc-123" {
			t.Errorf("expected 'trace-abc-123', got %q", id)
		}
	})

	t.Run("returns empty string when not present", func(t *testing.T) {
		ctx := context.Background()
		id := GetTraceIDFromContext(ctx)
		if id != "" {
			t.Errorf("expected empty string, got %q", id)
		}
	})

	t.Run("returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), traceIDContextKey, 12345)
		id := GetTraceIDFromContext(ctx)
		if id != "" {
			t.Errorf("expected empty string for non-string value, got %q", id)
		}
	})
}
