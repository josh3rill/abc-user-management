package repositories

import (
    "abc-user-management/internal/models"
)

// UserRepository defines the contract for user data access
type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    FindAll(page, limit int, search string) ([]models.User, int64, error)
    Update(user *models.User) error
    Delete(id uint) error
}

// RepositoryFactory creates repository based on configuration
type RepositoryFactory struct {
    dbType string
    db     interface{}
}

func NewRepositoryFactory(dbType string, db interface{}) *RepositoryFactory {
    return &RepositoryFactory{
        dbType: dbType,
        db:     db,
    }
}

// GetUserRepository returns appropriate repository implementation
func (rf *RepositoryFactory) GetUserRepository() UserRepository {
    switch rf.dbType {
    case "mysql":
        return NewMySQLUserRepository(rf.db)
    case "mongodb":
        return NewMongoUserRepository(rf.db)
    default:
        return NewMySQLUserRepository(rf.db)
    }
}