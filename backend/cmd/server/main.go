package main

import (
	"abc-user-management/internal/config"
	"abc-user-management/internal/events"
	"abc-user-management/internal/repositories"
	"abc-user-management/internal/routes"
	"abc-user-management/internal/services"
	"abc-user-management/pkg/cache"
	"abc-user-management/pkg/database"
	"abc-user-management/pkg/messaging"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	// Initialize configuration from environment
	cfg := config.Load()

	// Setup database connection with retry logic
	log.Println("Initializing database connection...")
	var db interface{}
	var err error

	// Retry database connection up to 10 times
	for i := 0; i < 10; i++ {
		db, err = database.InitializeDB(cfg)
		if err == nil {
			log.Println("Database connected successfully!")
			break
		}
		log.Printf("Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after multiple attempts: %v", err)
	}

	// Initialize Redis cache for performance optimization
	log.Println("Initializing Redis cache...")
	var redisClient *cache.RedisClient
	redisClient = cache.InitRedis(cfg)
	if redisClient != nil {
		// Test Redis connection
		err := redisClient.Set("test", "connection", 1*time.Second)
		if err != nil {
			log.Printf("Warning: Redis connection failed: %v. Running without cache.", err)
			redisClient = nil
		} else {
			log.Println("Redis cache connected successfully!")
			redisClient.Delete("test")
		}
	}

	// Setup message broker for event-driven architecture
	log.Println("Initializing RabbitMQ...")
	var msgBroker interface{}
	msgBroker, err = messaging.InitRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Warning: Message broker unavailable: %v. Running without events.", err)
		msgBroker = nil
	} else {
		log.Println("RabbitMQ connected successfully!")
		// Start event consumers in background
		go func() {
			if conn, ok := msgBroker.(*amqp.Connection); ok {
				if err := events.StartConsumers(conn, db); err != nil {
					log.Printf("Failed to start event consumers: %v", err)
				}
			}
		}()
	}

	// Create default admin user if not exists
	log.Println("Ensuring default admin user exists...")
	ensureDefaultAdmin(cfg, db)

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router with recovery middleware
	router := gin.New()
	router.Use(gin.Recovery())

	// Setup all routes and middleware
	routes.SetupRoutes(router, db, redisClient, msgBroker)

	// Start server on configured port
	log.Printf("🚀 ABC User Management API starting on port %s", cfg.Port)
	log.Printf("📝 Environment: %s", cfg.Environment)
	log.Printf("💾 Database: %s", cfg.DBType)
	log.Printf("📧 Default Admin: %s / %s", cfg.AdminEmail, cfg.AdminPassword)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ensureDefaultAdmin creates the default admin user if it doesn't exist
func ensureDefaultAdmin(cfg *config.Config, db interface{}) {
	// Create repository
	repoFactory := repositories.NewRepositoryFactory(cfg.DBType, db)
	userRepo := repoFactory.GetUserRepository()

	// Create auth service
	authService := services.NewAuthService(userRepo)

	// Try to create default admin
	err := authService.CreateDefaultAdmin(cfg.AdminEmail, cfg.AdminPassword)
	if err != nil {
		log.Printf("Default admin may already exist or creation failed: %v", err)
	} else {
		log.Println("Default admin user created successfully!")
	}
}
