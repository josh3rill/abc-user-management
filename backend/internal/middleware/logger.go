package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

// RequestLogger logs all incoming HTTP requests
func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer for request duration
        startTime := time.Now()
        
        // Process request
        c.Next()
        
        // Calculate request duration
        duration := time.Since(startTime)
        
        // Log request details after processing
        logger.WithFields(logrus.Fields{
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "status":     c.Writer.Status(),
            "duration":   duration,
            "ip":         c.ClientIP(),
            "user_agent": c.Request.UserAgent(),
            "error":      c.Errors.String(),
        }).Info("Request processed")
    }
}