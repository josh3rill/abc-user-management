package utils

import (
    "errors"
    "regexp"
    "strings"
)

var (
    ErrUserTooYoung  = errors.New("user must be 18 years or older")
    ErrEmailRequired = errors.New("email is required")
    ErrInvalidEmail  = errors.New("invalid email format")
    ErrPasswordWeak  = errors.New("password is too weak")
    ErrNameTooShort  = errors.New("name must be at least 2 characters")
    ErrNameTooLong   = errors.New("name must not exceed 100 characters")
)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
    // Basic email regex pattern
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return emailRegex.MatchString(strings.ToLower(email))
}

// ValidateAge checks if user meets minimum age requirement
func ValidateAge(age int) bool {
    return age >= 18
}

// ValidateName checks if name meets requirements
func ValidateName(name string) error {
    // Trim spaces and check length
    name = strings.TrimSpace(name)
    
    if len(name) < 2 {
        return ErrNameTooShort
    }
    
    if len(name) > 100 {
        return ErrNameTooLong
    }
    
    // Check for valid characters (letters, spaces, hyphens, apostrophes)
    validNameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
    if !validNameRegex.MatchString(name) {
        return errors.New("name contains invalid characters")
    }
    
    return nil
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
    if len(password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    
    // Optional: Check for password complexity
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    
    if !hasUpper || !hasLower || !hasNumber {
        return ErrPasswordWeak
    }
    
    return nil
}

// SanitizeInput removes potentially harmful characters
func SanitizeInput(input string) string {
    // Remove HTML tags and script injections
    htmlRegex := regexp.MustCompile(`<[^>]*>`)
    sanitized := htmlRegex.ReplaceAllString(input, "")
    
    // Trim excessive whitespace
    sanitized = strings.TrimSpace(sanitized)
    
    return sanitized
}

// ValidateUserInput performs all user input validations
func ValidateUserInput(name, email string, age int) error {
    // Validate name
    if err := ValidateName(name); err != nil {
        return err
    }
    
    // Validate email
    if !ValidateEmail(email) {
        return ErrInvalidEmail
    }
    
    // Validate age
    if !ValidateAge(age) {
        return ErrUserTooYoung
    }
    
    return nil
}