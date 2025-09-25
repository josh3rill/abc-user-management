package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"abc-user-management/internal/controllers"
	"abc-user-management/internal/middleware"
	"abc-user-management/internal/models"
	"abc-user-management/internal/services"
	"abc-user-management/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
	users  map[uint]*models.User
	nextID uint
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[uint]*models.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	if args.Error(0) == nil {
		user.ID = m.nextID
		m.users[m.nextID] = user
		m.nextID++
	}
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) FindAll(page, limit int, search string) ([]models.User, int64, error) {
	args := m.Called(page, limit, search)

	var users []models.User
	for _, user := range m.users {
		if search == "" ||
			(search != "" && (contains(user.Name, search) || contains(user.Email, search))) {
			users = append(users, *user)
		}
	}

	// Apply pagination
	start := (page - 1) * limit
	end := start + limit
	if start > len(users) {
		return []models.User{}, int64(len(users)), nil
	}
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], int64(len(users)), args.Error(2)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	if args.Error(0) == nil {
		m.users[user.ID] = user
	}
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	if args.Error(0) == nil {
		delete(m.users, id)
	}
	return args.Error(0)
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}

// APITestSuite for integration tests
type APITestSuite struct {
	suite.Suite
	router   *gin.Engine
	mockRepo *MockUserRepository
	token    string
}

// SetupSuite runs once before all tests
func (suite *APITestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// SetupTest runs before each test
func (suite *APITestSuite) SetupTest() {
	// Create mock repository
	suite.mockRepo = NewMockUserRepository()

	// Setup mock expectations for all methods
	suite.mockRepo.On("Create", mock.Anything).Return(nil)
	suite.mockRepo.On("FindByID", mock.Anything).Return(nil, nil)
	suite.mockRepo.On("FindByEmail", mock.Anything).Return(nil, nil)
	suite.mockRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil)
	suite.mockRepo.On("Update", mock.Anything).Return(nil)
	suite.mockRepo.On("Delete", mock.Anything).Return(nil)

	// Create test admin user
	hashedPassword, _ := utils.HashPassword("Admin@123456")
	adminUser := &models.User{
		ID:       1,
		Name:     "Admin User",
		Email:    "admin@abc.com",
		Password: hashedPassword,
		Age:      30,
		Role:     "admin",
		Active:   true,
	}
	suite.mockRepo.users[1] = adminUser
	suite.mockRepo.nextID = 2

	// Generate admin token for tests
	suite.token, _ = utils.GenerateJWT(adminUser.ID, adminUser.Email, adminUser.Role)

	// Setup router with mock services
	suite.setupRouter()
}

func (suite *APITestSuite) setupRouter() {
	suite.router = gin.New()

	// Add middleware
	suite.router.Use(middleware.CORS())

	// Initialize services with mock repository
	userService := services.NewUserService(suite.mockRepo, nil, nil) // No cache/events for testing
	authService := services.NewAuthService(suite.mockRepo)

	// Initialize controllers
	logger := logrus.New()
	userController := controllers.NewUserController(userService, logger)
	authController := controllers.NewAuthController(authService, logger)

	// Setup routes
	v1 := suite.router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/login", authController.Login)
		v1.GET("/default-credentials", authController.GetDefaultCredentials)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/users", userController.GetUsers)
			protected.GET("/users/:id", userController.GetUser)
			protected.POST("/users", userController.CreateUser)
			protected.PUT("/users/:id", userController.UpdateUser)
			protected.DELETE("/users/:id", userController.DeleteUser)
		}
	}

	// Health check
	suite.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
}

// Test Health Check
func (suite *APITestSuite) TestHealthCheck() {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "healthy", response["status"])
}

// Test Login Success
func (suite *APITestSuite) TestLoginSuccess() {
	loginData := map[string]string{
		"email":    "admin@abc.com",
		"password": "Admin@123456",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(suite.T(), response, "token")
	assert.Contains(suite.T(), response, "user")
}

// Test Login Failure
func (suite *APITestSuite) TestLoginFailure() {
	loginData := map[string]string{
		"email":    "wrong@email.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// Test Create User Success
func (suite *APITestSuite) TestCreateUserSuccess() {
	userData := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@test.com",
		"age":      25,
		"password": "Password123",
		"role":     "user",
		"active":   true,
	}
	body, _ := json.Marshal(userData)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(suite.T(), response, "data")
}

// Test Create User - Age Validation
func (suite *APITestSuite) TestCreateUserAgeValidation() {
	userData := map[string]interface{}{
		"name":     "Young User",
		"email":    "young@test.com",
		"age":      16, // Under 18
		"password": "Password123",
	}
	body, _ := json.Marshal(userData)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(suite.T(), response["error"], "18")
}

// Test Create User - Duplicate Email
func (suite *APITestSuite) TestCreateUserDuplicateEmail() {
	// First user
	userData1 := map[string]interface{}{
		"name":     "User One",
		"email":    "duplicate@test.com",
		"age":      25,
		"password": "Password123",
	}
	body1, _ := json.Marshal(userData1)

	req1 := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+suite.token)
	w1 := httptest.NewRecorder()
	suite.router.ServeHTTP(w1, req1)

	// Second user with same email
	userData2 := map[string]interface{}{
		"name":     "User Two",
		"email":    "duplicate@test.com",
		"age":      26,
		"password": "Password456",
	}
	body2, _ := json.Marshal(userData2)

	req2 := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+suite.token)
	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	assert.Equal(suite.T(), http.StatusConflict, w2.Code)
}

// Test Get Users List
func (suite *APITestSuite) TestGetUsersList() {
	// Create test users
	testUsers := []map[string]interface{}{
		{"name": "User 1", "email": "user1@test.com", "age": 25, "password": "Pass123"},
		{"name": "User 2", "email": "user2@test.com", "age": 30, "password": "Pass123"},
		{"name": "User 3", "email": "user3@test.com", "age": 35, "password": "Pass123"},
	}

	for _, userData := range testUsers {
		body, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.token)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
	}

	// Get users list
	req := httptest.NewRequest("GET", "/api/v1/users?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(suite.T(), response, "data")
	assert.Contains(suite.T(), response, "pagination")
}

// Test Pagination
func (suite *APITestSuite) TestPagination() {
	// Create 10 test users
	for i := 1; i <= 10; i++ {
		userData := map[string]interface{}{
			"name":     fmt.Sprintf("Test User %d", i),
			"email":    fmt.Sprintf("test%d@example.com", i),
			"age":      20 + i,
			"password": "Password123",
		}
		body, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.token)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
	}

	// Test page 1 with limit 5
	req := httptest.NewRequest("GET", "/api/v1/users?page=1&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	pagination := response["pagination"].(map[string]interface{})
	assert.Equal(suite.T(), float64(1), pagination["page"])
	assert.Equal(suite.T(), float64(5), pagination["limit"])
}

// Test Search Functionality
func (suite *APITestSuite) TestSearchUsers() {
	// Create users with different names
	users := []map[string]interface{}{
		{"name": "John Smith", "email": "john@test.com", "age": 25, "password": "Pass123"},
		{"name": "Jane Doe", "email": "jane@test.com", "age": 30, "password": "Pass123"},
		{"name": "John Doe", "email": "johndoe@test.com", "age": 35, "password": "Pass123"},
	}

	for _, userData := range users {
		body, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+suite.token)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
	}

	// Search for "John"
	req := httptest.NewRequest("GET", "/api/v1/users?search=John", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	users = response["data"].([]map[string]interface{})
	// Should find users with "John" in name
	for _, user := range users {
		assert.Contains(suite.T(), user["name"], "John")
	}
}

// Test Get Single User
func (suite *APITestSuite) TestGetSingleUser() {
	// Create a user first
	userData := map[string]interface{}{
		"name":     "Test User",
		"email":    "getuser@test.com",
		"age":      25,
		"password": "Password123",
	}
	body, _ := json.Marshal(userData)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	userID := createResponse["data"].(map[string]interface{})["id"]

	// Get the user
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(suite.T(), response, "data")
}

// Test Update User
func (suite *APITestSuite) TestUpdateUser() {
	// Create a user
	userData := map[string]interface{}{
		"name":     "Original Name",
		"email":    "update@test.com",
		"age":      25,
		"password": "Password123",
	}
	body, _ := json.Marshal(userData)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	userID := createResponse["data"].(map[string]interface{})["id"]

	// Update the user
	updateData := map[string]interface{}{
		"name": "Updated Name",
		"age":  30,
	}
	body, _ = json.Marshal(updateData)

	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%v", userID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Test Delete User
func (suite *APITestSuite) TestDeleteUser() {
	// Create a user
	userData := map[string]interface{}{
		"name":     "Delete Me",
		"email":    "delete@test.com",
		"age":      25,
		"password": "Password123",
	}
	body, _ := json.Marshal(userData)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	userID := createResponse["data"].(map[string]interface{})["id"]

	// Delete the user
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Verify user is deleted
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

// Test Authorization - No Token
func (suite *APITestSuite) TestAuthorizationNoToken() {
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	// No Authorization header
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// Test Authorization - Invalid Token
func (suite *APITestSuite) TestAuthorizationInvalidToken() {
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// Test Invalid JSON
func (suite *APITestSuite) TestInvalidJSON() {
	invalidJSON := `{"name": "Test", "age": }`

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// Test Empty Search
func (suite *APITestSuite) TestEmptySearch() {
	// Get all users with empty search
	req := httptest.NewRequest("GET", "/api/v1/users?search=", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Run the test suite
func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
