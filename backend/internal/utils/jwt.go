package utils

import (
    "errors"
    "time"
    "github.com/dgrijalva/jwt-go"
)

// JWT secret key (should be in environment variable)
var jwtSecret = []byte("your-super-secret-jwt-key-change-in-production")

// Claims represents JWT claims structure
type Claims struct {
    UserID    uint   `json:"user_id"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    jwt.StandardClaims
}

// GenerateJWT creates a new JWT token for authenticated user
func GenerateJWT(userID uint, email, role string) (string, error) {
    // Set token expiration time (24 hours)
    expirationTime := time.Now().Add(24 * time.Hour)
    
    // Create claims with user info
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
            IssuedAt:  time.Now().Unix(),
            Issuer:    "abc-user-management",
        },
    }
    
    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign token with secret key
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    
    return tokenString, nil
}

// ValidateJWT verifies and parses JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    // Parse token with claims
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Check if token is valid
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    // Check if token is expired
    if time.Now().Unix() > claims.ExpiresAt {
        return nil, errors.New("token expired")
    }
    
    return claims, nil
}

// RefreshJWT creates a new token with extended expiration
func RefreshJWT(oldToken string) (string, error) {
    // Validate old token
    claims, err := ValidateJWT(oldToken)
    if err != nil {
        return "", err
    }
    
    // Generate new token with same claims but new expiration
    return GenerateJWT(claims.UserID, claims.Email, claims.Role)
}

// SetJWTSecret updates the JWT secret (for testing purposes)
func SetJWTSecret(secret string) {
    jwtSecret = []byte(secret)
}