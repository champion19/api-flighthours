package middleware

import (
	"sync"
	"time"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter wraps a rate.Limiter with a last-seen timestamp for cleanup.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter implements per-IP token bucket rate limiting.
type RateLimiter struct {
	ips   sync.Map // map[string]*ipLimiter
	rps   rate.Limit
	burst int
}

// NewRateLimiter creates a new rate limiter with the given requests-per-second
// and burst size. It starts a background goroutine that cleans up stale entries
// every 3 minutes.
//
// Example: NewRateLimiter(0.2, 5) → 1 request every 5 seconds, burst of 5.
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		rps:   rate.Limit(rps),
		burst: burst,
	}

	go rl.cleanup()
	return rl
}

// cleanup removes IP entries that haven't been seen for 3+ minutes.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.ips.Range(func(key, value any) bool {
			entry := value.(*ipLimiter)
			if time.Since(entry.lastSeen) > 3*time.Minute {
				rl.ips.Delete(key)
			}
			return true
		})
	}
}

// getLimiter returns the rate.Limiter for the given IP, creating one if needed.
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	if v, ok := rl.ips.Load(ip); ok {
		entry := v.(*ipLimiter)
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(rl.rps, rl.burst)
	entry := &ipLimiter{limiter: limiter, lastSeen: time.Now()}
	rl.ips.Store(ip, entry)
	return limiter
}

// Limit returns a Gin middleware that enforces the rate limit.
// If the limit is exceeded, it sets ErrRateLimitExceeded and aborts.
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			traceID := GetRequestID(c)
			log := log.WithTraceID(traceID)
			log.Warn(logger.LogMiddlewareRateLimitHit,
				"client_ip", ip,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)
			_ = c.Error(domain.ErrRateLimitExceeded)
			c.Abort()
			return
		}

		c.Next()
	}
}
