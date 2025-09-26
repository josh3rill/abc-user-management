package controllers

import (
	"abc-user-management/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	authService *services.AuthService
	logger      *logrus.Logger
}

// NewAuthController creates authentication controller
func NewAuthController(authService *services.AuthService, logger *logrus.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		logger:      logger,
	}
}

// @Summary User Login
// @Description Authenticates user and returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "User Credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router// Login handles user authentication requests
func (c *AuthController) Login(ctx *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validate input credentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid credentials format",
		})
		return
	}

	// Authenticate through service
	token, user, err := c.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		c.logger.WithError(err).WithField("email", credentials.Email).Warn("Login attempt failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Log successful authentication
	c.logger.WithField("user_id", user.ID).Info("User logged in successfully")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// @Summary Get Default Admin Credentials
// @Description Provides default admin login details for demo/testing
// @Tags Auth
// @Accept json
// @Produce json
// @Parameters
// @Success 200
// @Failure 404
// @Router// GetDefaultCredentials provides default admin credentials for testing
// GetDefaultCredentials returns admin login details for demo
func (c *AuthController) GetDefaultCredentials(ctx *gin.Context) {
	// This endpoint provides default admin credentials for easy testing
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Default admin credentials for testing",
		"credentials": gin.H{
			"email":    "admin@abc.com",
			"password": "Admin@123456",
		},
		"note": "Click to copy these credentials on the login page",
	})
}
