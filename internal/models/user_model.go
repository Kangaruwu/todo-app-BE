package models

import (
	"time"

	"github.com/google/uuid"
)

// User model represents a user account

type UserRoleEnum string

const (
	AdminRole UserRoleEnum = "admin"
	UserRole  UserRoleEnum = "user"
)

type emailValidationStatusEnum string

const (
	UnconfirmedStatus emailValidationStatusEnum = "unconfirmed"
	ConfirmedStatus   emailValidationStatusEnum = "confirmed"
	PendingStatus     emailValidationStatusEnum = "pending"
)

type UserAccount struct {
	UserID   uuid.UUID    `json:"user_id" db:"user_id"`
	UserRole UserRoleEnum `json:"user_role" db:"user_role"`
}

type UserLoginData struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Username string    `json:"user_name" db:"user_name"`

	PasswordHash  string `json:"-" db:"password_hash"`  // not returned in JSON
	HashAlgorithm string `json:"-" db:"hash_algorithm"` // not returned in JSON

	EmailAddress string `json:"email_address" db:"email_address"`

	ConfirmationToken     string                    `json:"confirmation_token" db:"confirmation_token"`
	ConfirmationTokenTime time.Time                 `json:"confirmation_token_time" db:"confirmation_token_time"`
	EmailValidationStatus emailValidationStatusEnum `json:"email_validation_status" db:"email_validation_status"`

	PasswordRecoveryToken string    `json:"password_recovery_token" db:"password_recovery_token"`
	RecoveryTokenTime     time.Time `json:"recovery_token_time" db:"recovery_token_time"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest struct represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest struct represents the request to update an existing user
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
}

// UserProfile represents user profile information
type UserProfile struct {
	UserID    uuid.UUID                 `json:"user_id" db:"user_id"`
	Username  string                    `json:"username" db:"user_name"`
	Email     string                    `json:"email" db:"email_address"`
	Role      UserRoleEnum              `json:"role" db:"user_role"`
	Status    emailValidationStatusEnum `json:"email_status" db:"email_validation_status"`
	CreatedAt time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at" db:"updated_at"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe_updated"`
	Email    string `json:"email" validate:"required,email" example:"john.updated@example.com"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users  []UserProfile `json:"users"`
	Total  int64         `json:"total"`
	Page   int           `json:"page"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalUsers         int64 `json:"total_users"`
	ActiveUsers        int64 `json:"active_users"`
	PendingUsers       int64 `json:"pending_users"`
	RegisteredToday    int64 `json:"registered_today"`
	RegisteredThisWeek int64 `json:"registered_this_week"`
}
