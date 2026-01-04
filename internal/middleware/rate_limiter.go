package middleware

import (
	"sync"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a token bucket rate limiter per IP address
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

type visitor struct {
	lastSeen time.Time
	tokens   int
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rate: maximum number of requests
// window: time window (e.g., 1 minute)
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// Cleanup old visitors every minute
	go rl.cleanupVisitors()

	return rl
}

// Limit returns a middleware that limits requests per IP
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.allow(ip) {
			response.TooManyRequests(c, "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// allow checks if the request from the given IP is allowed
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: time.Now(),
			tokens:   rl.rate - 1,
		}
		rl.mu.Unlock()
		return true
	}
	rl.mu.Unlock()

	v.mu.Lock()
	defer v.mu.Unlock()

	// Refill tokens based on elapsed time
	elapsed := time.Since(v.lastSeen)
	if elapsed >= rl.window {
		v.tokens = rl.rate
		v.lastSeen = time.Now()
	}

	if v.tokens > 0 {
		v.tokens--
		v.lastSeen = time.Now()
		return true
	}

	return false
}

// cleanupVisitors removes visitors that haven't been seen in a while
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			v.mu.Lock()
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
			v.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// Common rate limiter presets

// RateLimitStrict limits to 10 requests per minute
func RateLimitStrict() gin.HandlerFunc {
	limiter := NewRateLimiter(10, 1*time.Minute)
	return limiter.Limit()
}

// RateLimitModerate limits to 60 requests per minute
func RateLimitModerate() gin.HandlerFunc {
	limiter := NewRateLimiter(60, 1*time.Minute)
	return limiter.Limit()
}

// RateLimitGenerous limits to 100 requests per minute
func RateLimitGenerous() gin.HandlerFunc {
	limiter := NewRateLimiter(100, 1*time.Minute)
	return limiter.Limit()
}
