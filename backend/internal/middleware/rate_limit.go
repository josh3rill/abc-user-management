package middleware

import (
    "net/http"
    "sync"
    "time"
    
    "github.com/gin-gonic/gin"
)

// RateLimiter tracks request rates per client
type RateLimiter struct {
    visitors map[string]*visitor
    mu       sync.RWMutex
    rate     int           // requests per duration
    duration time.Duration // time window
}

// visitor tracks request count for a client
type visitor struct {
    count     int
    lastReset time.Time
}

// NewRateLimiter creates rate limiting middleware
func NewRateLimiter(rate int, duration time.Duration) *RateLimiter {
    // Initialize limiter with cleanup routine
    rl := &RateLimiter{
        visitors: make(map[string]*visitor),
        rate:     rate,
        duration: duration,
    }
    
    // Start cleanup goroutine for expired entries
    go rl.cleanupVisitors()
    
    return rl
}

// Middleware returns Gin middleware function
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get client identifier (IP address)
        clientIP := c.ClientIP()
        
        // Check if request exceeds rate limit
        if !rl.allowRequest(clientIP) {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded. Please try again later.",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// allowRequest checks if client can make request
func (rl *RateLimiter) allowRequest(clientIP string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    // Get or create visitor record
    v, exists := rl.visitors[clientIP]
    now := time.Now()
    
    if !exists {
        // First request from this client
        rl.visitors[clientIP] = &visitor{
            count:     1,
            lastReset: now,
        }
        return true
    }
    
    // Check if window has expired
    if now.Sub(v.lastReset) > rl.duration {
        // Reset counter for new window
        v.count = 1
        v.lastReset = now
        return true
    }
    
    // Check if limit exceeded
    if v.count >= rl.rate {
        return false
    }
    
    // Increment counter
    v.count++
    return true
}

// cleanupVisitors removes expired visitor entries
func (rl *RateLimiter) cleanupVisitors() {
    for {
        time.Sleep(rl.duration)
        
        rl.mu.Lock()
        now := time.Now()
        
        // Remove visitors with expired windows
        for ip, v := range rl.visitors {
            if now.Sub(v.lastReset) > rl.duration*2 {
                delete(rl.visitors, ip)
            }
        }
        
        rl.mu.Unlock()
    }
}