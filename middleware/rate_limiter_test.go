package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRateLimiterRouter(rl *RateLimiter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(rl.Limit())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func TestRateLimiter_AllowsNormalTraffic(t *testing.T) {
	rl := NewRateLimiter(10, 5) // generous: 10 rps, burst 5

	r := setupRateLimiterRouter(rl)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected status 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_BlocksExcessiveRequests(t *testing.T) {
	rl := NewRateLimiter(0.001, 1) // very strict: burst 1

	handlerCalled := false

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(rl.Limit())
	r.POST("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First request passes (within burst)
	handlerCalled = false
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !handlerCalled {
		t.Fatal("first request should reach the handler")
	}

	// Second request should be blocked
	handlerCalled = false
	req = httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if handlerCalled {
		t.Fatal("second request should have been blocked by rate limiter")
	}
}

func TestRateLimiter_DifferentIPsIndependent(t *testing.T) {
	rl := NewRateLimiter(0.001, 1) // burst 1 per IP

	handlerReached := make(map[string]bool)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(rl.Limit())
	r.POST("/test", func(c *gin.Context) {
		handlerReached[c.ClientIP()] = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// IP-A first request — should pass
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "1.1.1.1:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !handlerReached["1.1.1.1"] {
		t.Fatal("IP-A first request should pass")
	}

	// IP-A second request — should be blocked
	handlerReached["1.1.1.1"] = false
	req = httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "1.1.1.1:12345"
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if handlerReached["1.1.1.1"] {
		t.Fatal("IP-A second request should be blocked")
	}

	// IP-B first request — should pass (independent limiter)
	req = httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "2.2.2.2:12345"
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !handlerReached["2.2.2.2"] {
		t.Fatal("IP-B first request should pass (independent from IP-A)")
	}
}

func TestRateLimiter_RecoversAfterWindow(t *testing.T) {
	// Use a rate that allows 1 request per 100ms with burst 1
	rl := NewRateLimiter(10, 1) // 10 rps = 1 req per 100ms, burst 1

	handlerCalled := false

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(rl.Limit())
	r.POST("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First request passes
	handlerCalled = false
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "3.3.3.3:12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !handlerCalled {
		t.Fatal("first request should pass")
	}

	// Second request immediately — should be blocked
	handlerCalled = false
	req = httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "3.3.3.3:12345"
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if handlerCalled {
		t.Fatal("second request (immediate) should be blocked")
	}

	// Wait for token to refill
	time.Sleep(150 * time.Millisecond)

	// Third request after wait — should pass again
	handlerCalled = false
	req = httptest.NewRequest(http.MethodPost, "/test", nil)
	req.RemoteAddr = "3.3.3.3:12345"
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !handlerCalled {
		t.Fatal("third request (after wait) should pass")
	}
}
