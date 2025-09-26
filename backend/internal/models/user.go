package models


import (
    "time"
    "gorm.io/gorm"

)

// Common validation errors for User


// User represents the core user entity in our system
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"size:100;not null"`
    Email     string         `json:"email" gorm:"uniqueIndex;size:255;not null"`
    Age       int            `json:"age" gorm:"not null"`
    Password  string         `json:"-" gorm:"size:255;not null"` // Never expose password in JSON
    Role      string         `json:"role" gorm:"size:50;default:'user'"`
    Active    bool           `json:"active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support
}

// TableName specifies the table name for GORM
func (User) TableName() string {
    return "users"
}

// Validate performs business rule validation
func (u *User) Validate() error {
    if u.Age < 18 {
        return ErrUserTooYoung
    }
    if u.Email == "" {
        return ErrEmailRequired
    }
    return nil

}
