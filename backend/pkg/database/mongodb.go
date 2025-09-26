package database

import (
    "context"
    "time"
    "abc-user-management/internal/config"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// InitializeMongoDB establishes connection to MongoDB
func InitializeMongoDB(cfg *config.Config) (*mongo.Database, error) {
    // Set client options with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Parse MongoDB URI and create client
    clientOptions := options.Client().ApplyURI(cfg.MongoURI)
    
    // Set connection pool size for production
    clientOptions.SetMaxPoolSize(50)
    clientOptions.SetMinPoolSize(10)
    
    // Connect to MongoDB
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }
    
    // Verify connection with ping
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }
    
    // Get database from connection
    database := client.Database("abc_users")
    
    return database, nil
}