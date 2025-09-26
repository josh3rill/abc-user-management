package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword generates bcrypt hash from plain password
func HashPassword(password string) (string, error) {
    // Generate hash with cost of 10 (reasonable security/performance balance)
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// CheckPasswordHash compares plain password with hash
func CheckPasswordHash(password, hash string) bool {
    // Compare password with hash
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// GenerateStrongHash creates a stronger hash (for sensitive operations)
func GenerateStrongHash(password string) (string, error) {
    // Use higher cost factor for more security (slower but more secure)
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}