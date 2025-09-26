package controllers

import (
	"abc-user-management/internal/models"
	"abc-user-management/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	service *services.UserService
	logger  *logrus.Logger
}

// NewUserController initializes controller with dependencies
func NewUserController(service *services.UserService, logger *logrus.Logger) *UserController {
	return &UserController{
		service: service,
		logger:  logger,
	}
}

// @Summary Create a new user
// @Description Creates a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [get]
// CreateUser handles POST /api/users endpoint
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User

	// Parse and validate request body
	if err := ctx.ShouldBindJSON(&user); err != nil {
		c.logger.WithError(err).Error("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	// Process user creation through service layer
	if err := c.service.CreateUser(&user); err != nil {
		c.logger.WithError(err).Error("Failed to create user")

		// Determine appropriate error response
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		} else if err.Error() == "user must be 18 or older" {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Log successful creation
	c.logger.WithField("user_id", user.ID).Info("User created successfully")

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// @Summary Get paginated list of users
// @Description Retrieves users with pagination and optional search
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {array} models.User
// @Failure 404 {object} map[string]string
// @Router /api/users [get]
// GetUsers handles GET /api/users with pagination and search
func (c *UserController) GetUsers(ctx *gin.Context) {
	// Parse query parameters with defaults
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Fetch users through service
	users, total, err := c.service.GetUsers(page, limit, search)
	if err != nil {
		c.logger.WithError(err).Error("Failed to fetch users")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// @Summary Get user by ID
// @Description Retrieves a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router// GetUser handles GET /api/users/:id
func (c *UserController) GetUser(ctx *gin.Context) {
	// Extract and validate user ID from URL
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Retrieve user from service
	user, err := c.service.GetUser(uint(id))
	if err != nil {
		c.logger.WithError(err).WithField("user_id", id).Error("User not found")
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// @Summary Update user by ID
// @Description Updates user details by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router// UpdateUser handles PUT /api/users/:id
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// Parse user ID from URL
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Parse update data from request body
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid update data",
		})
		return
	}

	// Process update through service
	if err := c.service.UpdateUser(uint(id), updates); err != nil {
		c.logger.WithError(err).WithField("user_id", id).Error("Update failed")

		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "user must be 18 or older" {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.logger.WithField("user_id", id).Info("User updated successfully")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

// @Summary Delete user by ID
// @Description Deletes a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router
// DeleteUser handles DELETE /api/users/:id
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// Extract user ID from route parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Process deletion through service
	if err := c.service.DeleteUser(uint(id)); err != nil {
		c.logger.WithError(err).WithField("user_id", id).Error("Deletion failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.logger.WithField("user_id", id).Info("User deleted successfully")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
