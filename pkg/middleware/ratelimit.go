package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//Http Request Limit

type RateLimiter struct {
	ips    map[string][]time.Time
	mu     sync.Mutex
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:    make(map[string][]time.Time),
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) cleanupOldRequests() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, requests := range rl.ips {
		var validRequests []time.Time
		for _, requestTime := range requests {
			if now.Sub(requestTime) <= rl.window {
				validRequests = append(validRequests, requestTime)
			}
		}

		if len(validRequests) > 0 {
			rl.ips[ip] = validRequests
		} else {
			delete(rl.ips, ip)
		}
	}
}

// Check ip request
func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Clear history request
	var validRequests []time.Time
	for _, requestTime := range rl.ips[ip] {
		if now.Sub(requestTime) <= rl.window {
			validRequests = append(validRequests, requestTime)
		}
	}

	// Check if HTTP requests exceed the max limit
	if len(validRequests) >= rl.limit {
		rl.ips[ip] = validRequests
		return false
	}

	// save new http request
	rl.ips[ip] = append(validRequests, now)
	return true
}

// Limit rate middleware
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	// Clear history request
	go func() {
		ticker := time.NewTicker(window / 2)
		defer ticker.Stop()

		for range ticker.C {
			limiter.cleanupOldRequests()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.isAllowed(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many HTTP requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
