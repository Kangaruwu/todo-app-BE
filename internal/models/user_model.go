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
	PasswordSalt  string `json:"-" db:"password_salt"`  // not returned in JSON
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
