package middleware

import (
    "net/http"
    "strings"
    "abc-user-management/internal/utils"
    
    "github.com/gin-gonic/gin"
)

// JWTAuth validates JWT tokens for protected routes
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header missing",
            })
            c.Abort()
            return
        }
        
        // Parse Bearer token format
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization format",
            })
            c.Abort()
            return
        }
        
        // Validate JWT token
        claims, err := utils.ValidateJWT(tokenParts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }
        
        // Store user info in context for handlers
        c.Set("userID", claims.UserID)
        c.Set("userEmail", claims.Email)
        c.Set("userRole", claims.Role)
        
        c.Next()
    }
}