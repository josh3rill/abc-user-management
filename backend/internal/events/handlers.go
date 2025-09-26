package events

import (
    "log"
    "time"
    "fmt"
)

// EventHandler defines the interface for event handlers
type EventHandler interface {
    Handle(data map[string]interface{}) error
}

// UserCreatedHandler processes user creation events
type UserCreatedHandler struct {
    emailService EmailService
    analytics    AnalyticsService
}

// Handle processes user.created events
func (h *UserCreatedHandler) Handle(data map[string]interface{}) error {
    // Extract user information from event data
    userID, _ := data["id"].(float64)
    email, _ := data["email"].(string)
    
    // Send welcome email to new user
    log.Printf("Sending welcome email to: %s", email)
    if h.emailService != nil {
        h.emailService.SendWelcomeEmail(email, uint(userID))
    }
    
    // Update analytics dashboard
    if h.analytics != nil {
        h.analytics.TrackUserSignup(uint(userID), time.Now())
    }
    
    // Additional processing like creating default settings
    log.Printf("Setting up default preferences for user ID: %v", userID)
    
    return nil
}

// UserUpdatedHandler manages user update events
type UserUpdatedHandler struct {
    cacheService  CacheInvalidator
    searchService SearchIndexer
}

// Handle processes user.updated events
func (h *UserUpdatedHandler) Handle(data map[string]interface{}) error {
    // Extract update information
    userID, _ := data["id"].(float64)
    changes, _ := data["changes"].(map[string]interface{})
    
    // Invalidate all related caches
    log.Printf("Invalidating caches for user ID: %v", userID)
    if h.cacheService != nil {
        h.cacheService.InvalidateUserCache(uint(userID))
    }
    
    // Update search index if email or name changed
    if _, hasEmail := changes["email"]; hasEmail {
        log.Printf("Updating search index for user ID: %v", userID)
        if h.searchService != nil {
            h.searchService.ReindexUser(uint(userID))
        }
    }
    
    // Log significant changes for audit trail
    log.Printf("User %v updated fields: %v", userID, changes)
    
    return nil
}

// UserDeletedHandler handles user deletion events
type UserDeletedHandler struct {
    cleanupService DataCleanupService
    notifier       NotificationService
}

// Handle processes user.deleted events
func (h *UserDeletedHandler) Handle(data map[string]interface{}) error {
    // Extract user ID from deletion event
    userID, _ := data["id"].(float64)
    timestamp, _ := data["timestamp"].(string)
    
    // Clean up user-related data from other services
    log.Printf("Starting cleanup for deleted user ID: %v", userID)
    
    if h.cleanupService != nil {
        // Remove user's files from storage
        h.cleanupService.RemoveUserFiles(uint(userID))
        
        // Archive user's activity logs
        h.cleanupService.ArchiveUserLogs(uint(userID))
        
        // Remove from mailing lists
        h.cleanupService.UnsubscribeFromServices(uint(userID))
    }
    
    // Notify administrators about deletion
    if h.notifier != nil {
        message := fmt.Sprintf("User ID %v was deleted at %s", userID, timestamp)
        h.notifier.NotifyAdmins("USER_DELETION", message)
    }
    
    log.Printf("Cleanup completed for user ID: %v", userID)
    
    return nil
}

// EventHandlerRegistry manages all event handlers
type EventHandlerRegistry struct {
    handlers map[string]EventHandler
}

// NewEventHandlerRegistry creates a new handler registry
func NewEventHandlerRegistry() *EventHandlerRegistry {
    return &EventHandlerRegistry{
        handlers: make(map[string]EventHandler),
    }
}

// Register adds a handler for specific event type
func (r *EventHandlerRegistry) Register(eventType string, handler EventHandler) {
    r.handlers[eventType] = handler
    log.Printf("Registered handler for event: %s", eventType)
}

// GetHandler retrieves handler for event type
func (r *EventHandlerRegistry) GetHandler(eventType string) (EventHandler, bool) {
    handler, exists := r.handlers[eventType]
    return handler, exists
}

// ProcessEvent routes event to appropriate handler
func (r *EventHandlerRegistry) ProcessEvent(eventType string, data map[string]interface{}) error {
    handler, exists := r.GetHandler(eventType)
    if !exists {
        log.Printf("No handler registered for event: %s", eventType)
        return fmt.Errorf("no handler for event type: %s", eventType)
    }
    
    // Execute handler with error recovery
    defer func() {
        if err := recover(); err != nil {
            log.Printf("Handler panic for %s: %v", eventType, err)
        }
    }()
    
    return handler.Handle(data)
}

// Service interfaces for dependency injection
type EmailService interface {
    SendWelcomeEmail(email string, userID uint) error
}

type AnalyticsService interface {
    TrackUserSignup(userID uint, timestamp time.Time) error
}

type CacheInvalidator interface {
    InvalidateUserCache(userID uint) error
}

type SearchIndexer interface {
    ReindexUser(userID uint) error
}

type DataCleanupService interface {
    RemoveUserFiles(userID uint) error
    ArchiveUserLogs(userID uint) error
    UnsubscribeFromServices(userID uint) error
}

type NotificationService interface {
    NotifyAdmins(eventType, message string) error
}