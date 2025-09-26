package database

import (
    "fmt"
    "time"
    "abc-user-management/internal/config"
    "abc-user-management/internal/models"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// InitializeDB sets up database connection based on configuration
func InitializeDB(cfg *config.Config) (*gorm.DB, error) {
    var db *gorm.DB
    var err error
    
    if cfg.DBType == "mysql" {
        db, err = initMySQL(cfg)
    } else if cfg.DBType == "mongodb" {
        // Return MongoDB connection wrapped in interface
        return initMongoDBWrapper(cfg)
    } else {
        // Default to MySQL
        db, err = initMySQL(cfg)
    }
    
    if err != nil {
        return nil, err
    }
    
    // Run migrations
    if err = runMigrations(db); err != nil {
        return nil, fmt.Errorf("failed to run migrations: %v", err)
    }
    
    return db, nil
}

// initMySQL creates MySQL connection
func initMySQL(cfg *config.Config) (*gorm.DB, error) {
    // Build MySQL DSN (Data Source Name)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
    )
    
    // Configure GORM
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().Local()
        },
    }
    
    // Open database connection
    db, err := gorm.Open(mysql.Open(dsn), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
    }
    
    // Get underlying SQL database to configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // Set connection pool settings for production
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // Test connection
    if err = sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping MySQL: %v", err)
    }
    
    return db, nil
}

// runMigrations executes database migrations
func runMigrations(db *gorm.DB) error {
    // Auto migrate models
    err := db.AutoMigrate(
        &models.User{},
        &models.Admin{},
    )
    
    if err != nil {
        return err
    }
    
    // Create indexes for better query performance
    if err := createIndexes(db); err != nil {
        return err
    }
    
    return nil
}

// createIndexes creates database indexes for optimization
func createIndexes(db *gorm.DB) error {
    // Create index on email for faster lookups
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
        return err
    }
    
    // Create index on created_at for sorting
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)").Error; err != nil {
        return err
    }
    
    // Create composite index for search queries
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_search ON users(name, email)").Error; err != nil {
        return err
    }
    
    return nil
}

// initMongoDBWrapper wraps MongoDB for compatibility
func initMongoDBWrapper(cfg *config.Config) (*gorm.DB, error) {
    // This is a placeholder for MongoDB integration
    // In production, you would implement a proper wrapper
    // or use a different interface that supports both databases
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