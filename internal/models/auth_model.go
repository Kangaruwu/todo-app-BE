package models

import (
	"github.com/google/uuid"
)	

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"securepassword123"`
}

// RegisterResponse represents user registration response
type RegisterResponse struct {
	UserID  uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message string    `json:"message" example:"User registered successfully"`
	Success bool      `json:"success" example:"true"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"securepassword123"`
}

// LoginResponse represents user login response
type LoginResponse struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int       `json:"expires_in" example:"3600"`
	UserID       uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// RecoverPasswordRequest represents password recovery request
type RecoverPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

// RecoverPasswordResponse represents password recovery response
type RecoverPasswordResponse struct {
	Message string `json:"message" example:"Password recovery email sent successfully"`
	Success bool   `json:"success" example:"true"`
}

// ResetPasswordRequest represents password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required" example:"recovery_token_123"`
	NewPassword string `json:"new_password" validate:"required,min=8" example:"newsecurepassword123"`
}

// ResetPasswordResponse represents password reset response
type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successfully"`
	Success bool   `json:"success" example:"true"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Success bool   `json:"success" example:"false"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
	Success bool   `json:"success" example:"true"`
}
