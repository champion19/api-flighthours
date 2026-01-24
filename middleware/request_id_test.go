package middleware

import (
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
