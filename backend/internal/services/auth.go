package services

import (
    "errors"
    "abc-user-management/internal/models"
    "abc-user-management/internal/repositories"
    "abc-user-management/internal/utils"
)

type AuthService struct {
    userRepo repositories.UserRepository
}

// NewAuthService creates new authentication service
func NewAuthService(repo repositories.UserRepository) *AuthService {
    return &AuthService{
        userRepo: repo,
    }
}

// Login authenticates user and returns JWT token
func (s *AuthService) Login(email, password string) (string, *models.User, error) {
    // Find user by email
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    // Verify password against hash
    if !utils.CheckPasswordHash(password, user.Password) {
        return "", nil, errors.New("invalid credentials")
    }
    
    // Generate JWT token for authenticated user
    token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
    if err != nil {
        return "", nil, errors.New("failed to generate token")
    }
    
    return token, user, nil
}

// CreateDefaultAdmin ensures default admin exists
func (s *AuthService) CreateDefaultAdmin(email, password string) error {
    // Check if admin already exists
    existing, _ := s.userRepo.FindByEmail(email)
    if existing != nil {
        return nil // Admin already exists, skip creation
    }
    
    // Hash the default password
    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        return err
    }
    
    // Create admin user with default values
    admin := &models.User{
        Name:     "Admin User",
        Email:    email,
        Password: hashedPassword,
        Age:      30,
        Role:     "admin",
        Active:   true,
    }
    
    return s.userRepo.Create(admin)
}

// RefreshToken generates new token for existing user
func (s *AuthService) RefreshToken(userID uint) (string, error) {
    // Fetch user to get current details
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return "", errors.New("user not found")
    }
    
    // Generate new token with updated expiry
    token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
    if err != nil {
        return "", errors.New("failed to refresh token")
    }
    
    return token, nil
}