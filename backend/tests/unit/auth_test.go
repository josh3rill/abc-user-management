package unit

import (
    "testing"
    "abc-user-management/internal/models"
    "abc-user-management/internal/services"
    "abc-user-management/internal/utils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation for testing
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(page, limit int, search string) ([]models.User, int64, error) {
    args := m.Called(page, limit, search)
    return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Update(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}

// TestAuthServiceLogin tests the login functionality
func TestAuthServiceLogin(t *testing.T) {
    // Setup mock repository
    mockRepo := new(MockUserRepository)
    authService := services.NewAuthService(mockRepo)
    
    // Create test user with hashed password
    hashedPassword, _ := utils.HashPassword("password123")
    testUser := &models.User{
        ID:       1,
        Email:    "test@example.com",
        Password: hashedPassword,
        Role:     "user",
    }
    
    // Test successful login
    t.Run("Successful login returns token", func(t *testing.T) {
        mockRepo.On("FindByEmail", "test@example.com").Return(testUser, nil).Once()
        
        token, user, err := authService.Login("test@example.com", "password123")
        
        assert.NoError(t, err)
        assert.NotEmpty(t, token)
        assert.Equal(t, testUser.Email, user.Email)
        mockRepo.AssertExpectations(t)
    })
    
    // Test invalid email
    t.Run("Invalid email returns error", func(t *testing.T) {
        mockRepo.On("FindByEmail", "invalid@example.com").Return(nil, models.ErrUserNotFound).Once()
        
        token, user, err := authService.Login("invalid@example.com", "password123")
        
        assert.Error(t, err)
        assert.Empty(t, token)
        assert.Nil(t, user)
        assert.Equal(t, "invalid credentials", err.Error())
        mockRepo.AssertExpectations(t)
    })
    
    // Test wrong password
    t.Run("Wrong password returns error", func(t *testing.T) {
        mockRepo.On("FindByEmail", "test@example.com").Return(testUser, nil).Once()
        
        token, user, err := authService.Login("test@example.com", "wrongpassword")
        
        assert.Error(t, err)
        assert.Empty(t, token)
        assert.Nil(t, user)
        assert.Equal(t, "invalid credentials", err.Error())
        mockRepo.AssertExpectations(t)
    })
}

// TestJWTGeneration tests JWT token generation and validation
func TestJWTGeneration(t *testing.T) {
    // Test token generation
    t.Run("Generate valid JWT token", func(t *testing.T) {
        token, err := utils.GenerateJWT(1, "test@example.com", "admin")
        
        assert.NoError(t, err)
        assert.NotEmpty(t, token)
    })
    
    // Test token validation
    t.Run("Validate JWT token", func(t *testing.T) {
        // Generate token first
        token, _ := utils.GenerateJWT(1, "test@example.com", "admin")
        
        // Validate the generated token
        claims, err := utils.ValidateJWT(token)
        
        assert.NoError(t, err)
        assert.NotNil(t, claims)
        assert.Equal(t, uint(1), claims.UserID)
        assert.Equal(t, "test@example.com", claims.Email)
        assert.Equal(t, "admin", claims.Role)
    })
    
    // Test invalid token
    t.Run("Invalid token returns error", func(t *testing.T) {
        claims, err := utils.ValidateJWT("invalid.token.here")
        
        assert.Error(t, err)
        assert.Nil(t, claims)
    })
}

// TestPasswordHashing tests password hashing and verification
func TestPasswordHashing(t *testing.T) {
    password := "MySecurePassword123!"
    
    // Test hashing
    t.Run("Hash password successfully", func(t *testing.T) {
        hash, err := utils.HashPassword(password)
        
        assert.NoError(t, err)
        assert.NotEmpty(t, hash)
        assert.NotEqual(t, password, hash)
    })
    
    // Test verification
    t.Run("Verify correct password", func(t *testing.T) {
        hash, _ := utils.HashPassword(password)
        
        isValid := utils.CheckPasswordHash(password, hash)
        assert.True(t, isValid)
    })
    
    // Test wrong password verification
    t.Run("Reject incorrect password", func(t *testing.T) {
        hash, _ := utils.HashPassword(password)
        
        isValid := utils.CheckPasswordHash("WrongPassword", hash)
        assert.False(t, isValid)
    })
}

// TestCreateDefaultAdmin tests default admin creation
func TestCreateDefaultAdmin(t *testing.T) {
    mockRepo := new(MockUserRepository)
    authService := services.NewAuthService(mockRepo)
    
    // Test creating new admin when none exists
    t.Run("Create admin when not exists", func(t *testing.T) {
        mockRepo.On("FindByEmail", "admin@abc.com").Return(nil, models.ErrUserNotFound).Once()
        mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()
        
        err := authService.CreateDefaultAdmin("admin@abc.com", "Admin@123456")
        
        assert.NoError(t, err)
        mockRepo.AssertExpectations(t)
    })
    
    // Test when admin already exists
    t.Run("Skip when admin exists", func(t *testing.T) {
        existingAdmin := &models.User{
            ID:    1,
            Email: "admin@abc.com",
            Role:  "admin",
        }
        mockRepo.On("FindByEmail", "admin@abc.com").Return(existingAdmin, nil).Once()
        
        err := authService.CreateDefaultAdmin("admin@abc.com", "Admin@123456")
        
        assert.NoError(t, err)
        mockRepo.AssertExpectations(t)
    })
}