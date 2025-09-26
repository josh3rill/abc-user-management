package routes

import (
	"abc-user-management/internal/controllers"
	"abc-user-management/internal/events"
	"abc-user-management/internal/middleware"
	"abc-user-management/internal/repositories"
	"abc-user-management/internal/services"
	"abc-user-management/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all API endpoints and middleware
func SetupRoutes(router *gin.Engine, db interface{}, cacheClient interface{}, broker interface{}) {
	// Swagger endpoint

	// Initialize logger for request tracking
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Setup middleware stack
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger(logger))
	router.Use(gin.Recovery())

	// Initialize repository with database abstraction
	repoFactory := repositories.NewRepositoryFactory("mysql", db)
	userRepo := repoFactory.GetUserRepository()

	// Type assert cache client
	var redisClient *cache.RedisClient
	if cacheClient != nil {
		if rc, ok := cacheClient.(*cache.RedisClient); ok {
			redisClient = rc
		} else if rc, ok := cacheClient.(*redis.Client); ok {
			// If it's a raw redis client, wrap it
			redisClient = &cache.RedisClient{
				Client: rc,
			}
		}
	}

	// Type assert and setup event publisher
	var eventPublisher *events.EventPublisher
	if broker != nil {
		if conn, ok := broker.(*amqp.Connection); ok {
			// Create event publisher from connection
			publisher, err := events.NewEventPublisher(conn)
			if err != nil {
				logger.WithError(err).Warn("Failed to create event publisher")
			} else {
				eventPublisher = publisher
			}
		} else if ep, ok := broker.(*events.EventPublisher); ok {
			// Already an event publisher
			eventPublisher = ep
		}
	}

	// Setup services with dependencies
	userService := services.NewUserService(userRepo, redisClient, eventPublisher)
	authService := services.NewAuthService(userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService, logger)
	authController := controllers.NewAuthController(authService, logger)

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		// Public endpoints - no authentication required
		public := v1.Group("/")
		{
			public.POST("/login", authController.Login)
			public.GET("/default-credentials", authController.GetDefaultCredentials)
		}

		// Protected endpoints - require valid JWT
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			// User management endpoints
			protected.GET("/users", userController.GetUsers)
			protected.GET("/users/:id", userController.GetUser)
			protected.POST("/users", userController.CreateUser)
			protected.PUT("/users/:id", userController.UpdateUser)
			protected.DELETE("/users/:id", userController.DeleteUser)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Health check endpoint for monitoring
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "ABC User Management API",
		})
	})
}
