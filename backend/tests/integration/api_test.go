package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// TestCreateUserAPI tests user creation endpoint
func TestCreateUserAPI(t *testing.T) {
    // Setup test router
    router := setupTestRouter()
    
    // Test cases for user creation
    tests := []struct {
        name       string
        payload    map[string]interface{}
        wantStatus int
        wantError  bool
    }{
        {
            name: "Valid user creation succeeds",
            payload: map[string]interface{}{
                "name":     "John Doe",
                "email":    "john@test.com",
                "age":      25,
                "password": "password123",
            },
            wantStatus: http.StatusCreated,
            wantError:  false,
        },
        {
            name: "User under 18 rejected",
            payload: map[string]interface{}{
                "name":     "Young User",
                "email":    "young@test.com",
                "age":      16,
                "password": "password123",
            },
            wantStatus: http.StatusBadRequest,
            wantError:  true,
        },
        {
            name: "Duplicate email rejected",
            payload: map[string]interface{}{
                "name":     "Duplicate User",
                "email":    "existing@test.com",
                "age":      30,
                "password": "password123",
            },
            wantStatus: http.StatusConflict,
            wantError:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Prepare request body
            jsonBody, _ := json.Marshal(tt.payload)
            req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Authorization", "Bearer test-token")
            
            // Record response
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            
            // Assert status code
            assert.Equal(t, tt.wantStatus, w.Code)
            
            // Parse response
            var response map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &response)
            
            if tt.wantError {
                assert.Contains(t, response, "error")
            } else {
                assert.Contains(t, response, "data")
            }
        })
    }
}

// TestGetUsersWithPagination tests user listing with pagination
func TestGetUsersWithPagination(t *testing.T) {
    router := setupTestRouter()
    
    // Test pagination parameters
    req := httptest.NewRequest("GET", "/api/v1/users?page=1&limit=10&search=john", nil)
    req.Header.Set("Authorization", "Bearer test-token")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Verify pagination response structure
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    assert.Contains(t, response, "data")
    assert.Contains(t, response, "pagination")
    
    pagination := response["pagination"].(map[string]interface{})
    assert.Contains(t, pagination, "page")
    assert.Contains(t, pagination, "limit")
    assert.Contains(t, pagination, "total")
}

// TestAuthenticationRequired tests protected endpoints
func TestAuthenticationRequired(t *testing.T) {
    router := setupTestRouter()
    
    // Request without authentication header
    req := httptest.NewRequest("GET", "/api/v1/users", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    // Should return unauthorized
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// setupTestRouter creates router for testing
func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    // Add test routes and middleware
    // This would typically use your actual route setup
    // with mock services for testing
    
    return router
}