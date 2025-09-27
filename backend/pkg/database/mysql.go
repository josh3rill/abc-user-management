package database

import (
	"abc-user-management/internal/config"
	"abc-user-management/internal/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitializeDB sets up database connection based on configuration
func InitializeDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	log.Printf("Initializing database connection to %s:%s/%s as user %s",
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)

	switch cfg.DBType {
	case "mysql":
		db, err = initMySQL(cfg)
	case "mongodb":
		return initMongoDBWrapper(cfg)
	default:
		// Default to MySQL
		db, err = initMySQL(cfg)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Run migrations - simplified without problematic index creation
	if err = runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database connection and migrations completed successfully")
	return db, nil
}

// initMySQL creates MySQL connection with retry logic
func initMySQL(cfg *config.Config) (*gorm.DB, error) {
	// Build MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	log.Printf("Connecting to MySQL with DSN: %s:***@tcp(%s:%s)/%s",
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Less verbose logging
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// Retry connection logic
	var db *gorm.DB
	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to MySQL (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL after %d attempts: %v", maxRetries, err)
	}

	// Get underlying SQL database to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL database: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %v", err)
	}

	log.Println("MySQL connection established successfully")
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Simply run AutoMigrate - let GORM handle everything
	err := db.AutoMigrate(
		&models.User{},
		&models.Admin{},
	)

	if err != nil {
		return fmt.Errorf("failed to run auto migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// initMongoDBWrapper wraps MongoDB for compatibility
func initMongoDBWrapper(_ *config.Config) (*gorm.DB, error) {
	return nil, fmt.Errorf("MongoDB support requires additional configuration")
}

// CloseDB closes database connection gracefully
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
