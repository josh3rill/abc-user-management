package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Environment   string
	DBType        string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	MongoURI      string
	JWTSecret     string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RabbitMQURL   string
	AdminEmail    string
	AdminPassword string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8085"),
		Environment:   getEnv("ENV", "development"),
		DBType:        getEnv("DB_TYPE", "mysql"),
		DBHost:        getEnv("DB_HOST", "mysql"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "abc_user"),
		DBPassword:    getEnv("DB_PASSWORD", "v12345"),
		DBName:        getEnv("DB_NAME", "abc_users"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017/abc_users"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret-change-this"),
		RedisHost:     getEnv("REDIS_HOST", "redis"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RabbitMQURL:   getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@abc.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "Admin@123456"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
