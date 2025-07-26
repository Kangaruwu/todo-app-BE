package models

import (
	"github.com/google/uuid"
)

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"Securep@ssword123"`
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
	Password string `json:"password" validate:"required,min=8" example:"Securep@ssword123"`
}

// LoginResponse represents user login response
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
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
	NewPassword string `json:"new_password" validate:"required,min=8" example:"Newsecurep@ssword123"`
}

// ResetPasswordResponse represents password reset response
type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successfully"`
	Success bool   `json:"success" example:"true"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"Securep@ssword123"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100" example:"Newsecurep@ssword123"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword" example:"Newsecurep@ssword123"`
}

// ChangePasswordResponse represents password change response
type ChangePasswordResponse struct {
	Message string `json:"message" example:"Password changed successfully"`
	Success bool   `json:"success" example:"true"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}

// EmailConfirmationResponse represents email confirmation response
type EmailConfirmationResponse struct {
	Message string `json:"message" example:"Email confirmed successfully"`
	Success bool   `json:"success" example:"true"`
}

// RefreshAccessTokenRequest represents request to refresh access token
type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshAccessTokenResponse represents response for access token refresh
type RefreshAccessTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // refresh token rotated
}	

