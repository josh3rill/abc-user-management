package repositories

import (
    "abc-user-management/internal/models"
    "gorm.io/gorm"
)

type MySQLUserRepository struct {
    db *gorm.DB
}

// NewMySQLUserRepository creates a new MySQL repository instance
func NewMySQLUserRepository(db interface{}) UserRepository {
    return &MySQLUserRepository{db: db.(*gorm.DB)}
}

// Create inserts a new user into the database
func (r *MySQLUserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

// FindByID retrieves a user by their unique identifier
func (r *MySQLUserRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByEmail locates a user using their email address
func (r *MySQLUserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// FindAll retrieves paginated users with optional search
func (r *MySQLUserRepository) FindAll(page, limit int, search string) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    query := r.db.Model(&models.User{})
    
    // Apply search filter if provided
    if search != "" {
        query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
    }
    
    // Get total count for pagination
    query.Count(&total)
    
    // Apply pagination
    offset := (page - 1) * limit
    err := query.Offset(offset).Limit(limit).Find(&users).Error
    
    return users, total, err
}

// Update modifies an existing user's information
func (r *MySQLUserRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

// Delete performs soft deletion of a user
func (r *MySQLUserRepository) Delete(id uint) error {
    return r.db.Delete(&models.User{}, id).Error
}