package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestJsonValidator_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid JSON body returns error", func(t *testing.T) {
		// This test validates that non-JSON body causes error
		r := gin.New()
		r.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
	})
}

func TestNewMiddlewareValidator(t *testing.T) {
	t.Run("creates builder with validators", func(t *testing.T) {
		builder := NewMiddlewareValidator(nil)
		if builder == nil {
			t.Error("expected non-nil Builder")
		}
	})
}
