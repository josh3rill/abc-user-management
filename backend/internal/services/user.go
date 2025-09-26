package services

import (
    "errors"
    "abc-user-management/internal/models"
    "abc-user-management/internal/repositories"
    "abc-user-management/internal/utils"
    "abc-user-management/pkg/cache"
    "abc-user-management/internal/events"
    "encoding/json"
    "fmt"
    "time"
)

type UserService struct {
    repo      repositories.UserRepository
    cache     *cache.RedisClient
    publisher *events.EventPublisher
}

// NewUserService creates a service with all dependencies injected
func NewUserService(repo repositories.UserRepository, cache *cache.RedisClient, pub *events.EventPublisher) *UserService {
    return &UserService{
        repo:      repo,
        cache:     cache,
        publisher: pub,
    }
}

// CreateUser handles user creation with validation and caching
func (s *UserService) CreateUser(user *models.User) error {
    // Validate business rules before persisting
    if err := user.Validate(); err != nil {
        return err
    }
    
    // Check for duplicate email
    existing, _ := s.repo.FindByEmail(user.Email)
    if existing != nil {
        return errors.New("email already exists")
    }
    
    // Hash password before storage
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        return err
    }
    user.Password = hashedPassword
    
    // Persist to database
    if err := s.repo.Create(user); err != nil {
        return err
    }
    
    // Invalidate cache for user list (if cache is available)
    if s.cache != nil {
        s.cache.Delete("users:list:*")
    }
    
    // Publish event for other services (if publisher is available)
    if s.publisher != nil {
        s.publisher.Publish("user.created", map[string]interface{}{
            "id": user.ID,
            "email": user.Email,
            "timestamp": time.Now(),
        })
    }
    
    return nil
}

// GetUser retrieves a user with caching strategy
func (s *UserService) GetUser(id uint) (*models.User, error) {
    // Try cache first for performance (if available)
    if s.cache != nil {
        cacheKey := fmt.Sprintf("user:%d", id)
        cached, err := s.cache.Get(cacheKey)
        if err == nil && cached != "" {
            var user models.User
            if json.Unmarshal([]byte(cached), &user) == nil {
                return &user, nil
            }
        }
    }
    
    // Cache miss or no cache, fetch from database
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache the result for future requests (if cache is available)
    if s.cache != nil {
        userData, _ := json.Marshal(user)
        cacheKey := fmt.Sprintf("user:%d", id)
        s.cache.Set(cacheKey, string(userData), 5*time.Minute)
    }
    
    return user, nil
}

// GetUsers returns paginated user list with search capability
func (s *UserService) GetUsers(page, limit int, search string) ([]models.User, int64, error) {
    // Generate cache key based on parameters
    cacheKey := fmt.Sprintf("users:list:%d:%d:%s", page, limit, search)
    
    // Check cache first (if available)
    if s.cache != nil {
        cached, err := s.cache.Get(cacheKey)
        if err == nil && cached != "" {
            var result struct {
                Users []models.User `json:"users"`
                Total int64        `json:"total"`
            }
            if json.Unmarshal([]byte(cached), &result) == nil {
                return result.Users, result.Total, nil
            }
        }
    }
    
    // Fetch from repository
    users, total, err := s.repo.FindAll(page, limit, search)
    if err != nil {
        return nil, 0, err
    }
    
    // Cache the results (if cache is available)
    if s.cache != nil {
        cacheData, _ := json.Marshal(map[string]interface{}{
            "users": users,
            "total": total,
        })
        s.cache.Set(cacheKey, string(cacheData), 2*time.Minute)
    }
    
    return users, total, nil
}

// UpdateUser modifies user data with validation
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return err
    }
    
    // Apply updates dynamically
    if name, ok := updates["name"].(string); ok {
        user.Name = name
    }
    if age, ok := updates["age"].(int); ok {
        user.Age = age
    }
    if email, ok := updates["email"].(string); ok {
        user.Email = email
    }
    
    // Validate after updates
    if err := user.Validate(); err != nil {
        return err
    }
    
    // Persist changes
    if err := s.repo.Update(user); err != nil {
        return err
    }
    
    // Invalidate relevant caches (if cache is available)
    if s.cache != nil {
        s.cache.Delete(fmt.Sprintf("user:%d", id))
        s.cache.Delete("users:list:*")
    }
    
    // Publish update event (if publisher is available)
    if s.publisher != nil {
        s.publisher.Publish("user.updated", map[string]interface{}{
            "id": id,
            "changes": updates,
            "timestamp": time.Now(),
        })
    }
    
    return nil
}

// DeleteUser removes a user and clears caches
func (s *UserService) DeleteUser(id uint) error {
    if err := s.repo.Delete(id); err != nil {
        return err
    }
    
    // Clear all related caches (if cache is available)
    if s.cache != nil {
        s.cache.Delete(fmt.Sprintf("user:%d", id))
        s.cache.Delete("users:list:*")
    }
    
    // Notify other services about deletion (if publisher is available)
    if s.publisher != nil {
        s.publisher.Publish("user.deleted", map[string]interface{}{
            "id": id,
            "timestamp": time.Now(),
        })
    }
    
    return nil
}