package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetErrorType(t *testing.T) {
	tests := []struct {
		status   int
		expected string
	}{
		{500, "server_error"},
		{502, "server_error"},
		{404, "not_found"},
		{401, "auth_error"},
		{403, "auth_error"},
		{400, "bad_request"},
		{422, "validation_error"},
		{409, "conflict"},
		{429, "rate_limit_error"},
		{418, "client_error"}, // I'm a teapot - other client error
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			result := getErrorType(tt.status)
			if result != tt.expected {
				t.Errorf("getErrorType(%d) = %q, expected %q", tt.status, result, tt.expected)
			}
		})
	}
}

func TestTrackMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("does not panic on request", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("TrackMetrics panicked: %v", r)
			}
		}()

		r := gin.New()
		r.Use(TrackMetrics())
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
	})

	t.Run("tracks error responses", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("TrackMetrics panicked: %v", r)
			}
		}()

		r := gin.New()
		r.Use(TrackMetrics())
		r.GET("/error", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "test"})
		})

		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
	})
}

func TestPrometheusMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrometheusMiddleware panicked: %v", r)
			}
		}()

		r := gin.New()
		r.Use(PrometheusMiddleware())
		r.GET("/metrics-test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/metrics-test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
	})
}

func TestRecordEmployeeRegistration(t *testing.T) {
	t.Run("does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RecordEmployeeRegistration panicked: %v", r)
			}
		}()
		RecordEmployeeRegistration()
	})
}

func TestRecordMessageCreated(t *testing.T) {
	t.Run("does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RecordMessageCreated panicked: %v", r)
			}
		}()
		RecordMessageCreated("test_module", "info")
	})
}
