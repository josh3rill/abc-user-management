package unit

import (
    "testing"
    "abc-user-management/internal/models"
)

// TestUserValidation verifies business rules are enforced
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        user    models.User
        wantErr bool
    }{
        {
            name: "Valid user passes validation",
            user: models.User{
                Name:  "John Doe",
                Email: "john@example.com",
                Age:   25,
            },
            wantErr: false,
        },
        {
            name: "User under 18 fails validation",
            user: models.User{
                Name:  "Young User",
                Email: "young@example.com",
                Age:   17,
            },
            wantErr: true,
        },
        {
            name: "User without email fails validation",
            user: models.User{
                Name: "No Email",
                Age:  20,
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// TestEmailUniqueness checks duplicate email handling
func TestEmailUniqueness(t *testing.T) {
    // This would require a mock repository
    // Implementation depends on your testing strategy
}