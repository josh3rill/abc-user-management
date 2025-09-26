package repositories

import (
    "context"
    "abc-user-management/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoUserRepository implements UserRepository for MongoDB
type MongoUserRepository struct {
    db         *mongo.Database
    collection *mongo.Collection
    ctx        context.Context
}

// NewMongoUserRepository creates MongoDB repository instance
func NewMongoUserRepository(database interface{}) UserRepository {
    db := database.(*mongo.Database)
    return &MongoUserRepository{
        db:         db,
        collection: db.Collection("users"),
        ctx:        context.Background(),
    }
}

// Create inserts new user document into MongoDB
func (r *MongoUserRepository) Create(user *models.User) error {
    // Convert to BSON and insert
    result, err := r.collection.InsertOne(r.ctx, user)
    if err != nil {
        return err
    }
    
    // Update user ID from inserted document
    if oid, ok := result.InsertedID.(uint); ok {
        user.ID = oid
    }
    
    return nil
}

// FindByID retrieves user document by ID
func (r *MongoUserRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    
    // Find document with matching ID
    filter := bson.M{"_id": id}
    err := r.collection.FindOne(r.ctx, filter).Decode(&user)
    
    if err == mongo.ErrNoDocuments {
        return nil, models.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// FindByEmail locates user by email in MongoDB
func (r *MongoUserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    
    // Create filter for email search
    filter := bson.M{"email": email}
    err := r.collection.FindOne(r.ctx, filter).Decode(&user)
    
    if err == mongo.ErrNoDocuments {
        return nil, models.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// FindAll retrieves paginated users with optional text search
func (r *MongoUserRepository) FindAll(page, limit int, search string) ([]models.User, int64, error) {
    var users []models.User
    
    // Build filter based on search criteria
    filter := bson.M{}
    if search != "" {
        // Text search on name and email fields
        filter = bson.M{
            "$or": []bson.M{
                {"name": bson.M{"$regex": search, "$options": "i"}},
                {"email": bson.M{"$regex": search, "$options": "i"}},
            },
        }
    }
    
    // Count total documents matching filter
    total, err := r.collection.CountDocuments(r.ctx, filter)
    if err != nil {
        return nil, 0, err
    }
    
    // Calculate pagination offset
    skip := int64((page - 1) * limit)
    
    // Set find options for pagination and sorting
    findOptions := options.Find()
    findOptions.SetLimit(int64(limit))
    findOptions.SetSkip(skip)
    findOptions.SetSort(bson.D{{"created_at", -1}}) // Sort by newest first
    
    // Execute query with pagination
    cursor, err := r.collection.Find(r.ctx, filter, findOptions)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(r.ctx)
    
    // Decode results into user slice
    for cursor.Next(r.ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        users = append(users, user)
    }
    
    return users, total, nil
}

// Update modifies existing user document
func (r *MongoUserRepository) Update(user *models.User) error {
    // Create filter to find document
    filter := bson.M{"_id": user.ID}
    
    // Prepare update document with all fields
    update := bson.M{
        "$set": bson.M{
            "name":       user.Name,
            "email":      user.Email,
            "age":        user.Age,
            "role":       user.Role,
            "active":     user.Active,
            "updated_at": user.UpdatedAt,
        },
    }
    
    // Execute update operation
    result, err := r.collection.UpdateOne(r.ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return models.ErrUserNotFound
    }
    
    return nil
}

// Delete performs soft delete on user document
func (r *MongoUserRepository) Delete(id uint) error {
    // Soft delete by setting deleted_at timestamp
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "deleted_at": bson.M{"$currentDate": true},
            "active":     false,
        },
    }
    
    result, err := r.collection.UpdateOne(r.ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return models.ErrUserNotFound
    }
    
    return nil
}

// CreateIndexes ensures MongoDB indexes are created
func (r *MongoUserRepository) CreateIndexes() error {
    // Create unique index on email field
    emailIndex := mongo.IndexModel{
        Keys:    bson.D{{"email", 1}},
        Options: options.Index().SetUnique(true),
    }
    
    // Create text index for search functionality
    textIndex := mongo.IndexModel{
        Keys: bson.D{
            {"name", "text"},
            {"email", "text"},
        },
    }
    
    // Create compound index for efficient querying
    queryIndex := mongo.IndexModel{
        Keys: bson.D{
            {"active", 1},
            {"created_at", -1},
        },
    }
    
    // Apply all indexes to collection
    _, err := r.collection.Indexes().CreateMany(r.ctx, []mongo.IndexModel{
        emailIndex,
        textIndex,
        queryIndex,
    })
    
    return err
}