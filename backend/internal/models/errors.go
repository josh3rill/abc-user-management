package models

import "errors"

var (
    ErrUserTooYoung   = errors.New("user is too young")
    ErrEmailRequired  = errors.New("email is required")
    ErrUserNotFound   = errors.New("user not found")
)
