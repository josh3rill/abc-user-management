package models

import (
    "time"
)

// Admin represents administrative user with extended privileges
type Admin struct {
    User
    Permissions   []string  `json:"permissions" gorm:"type:json"`
    LastLoginAt   time.Time `json:"last_login_at"`
    LoginAttempts int       `json:"login_attempts"`
    LockedUntil   time.Time `json:"locked_until"`
}

// HasPermission checks if admin has specific permission
func (a *Admin) HasPermission(permission string) bool {
    for _, p := range a.Permissions {
        if p == permission || p == "*" {
            return true
        }
    }
    return false
}

// IsLocked checks if admin account is temporarily locked
func (a *Admin) IsLocked() bool {
    return time.Now().Before(a.LockedUntil)
}

// IncrementLoginAttempts tracks failed login attempts
func (a *Admin) IncrementLoginAttempts() {
    a.LoginAttempts++
    
    // Lock account after 5 failed attempts
    if a.LoginAttempts >= 5 {
        a.LockedUntil = time.Now().Add(30 * time.Minute)
    }
}

// ResetLoginAttempts clears failed login counter
func (a *Admin) ResetLoginAttempts() {
    a.LoginAttempts = 0
    a.LastLoginAt = time.Now()
    a.LockedUntil = time.Time{}
}