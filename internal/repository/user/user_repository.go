package user_repository

import (
	"context"

	"go-backend-todo/internal/models"

	"github.com/google/uuid"
)

// UserRepository interface defines methods for interacting with user data
type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, req *models.RegisterRequest, verificationToken string) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error)
	GetByEmail(ctx context.Context, email string) (*models.UserProfile, error)
	GetByUsername(ctx context.Context, username string) (*models.UserProfile, error)
	Update(ctx context.Context, user *models.UserAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error

	// Token version operations
	IncrementTokenVersion(ctx context.Context, userID uuid.UUID) error
	GetTokenVersion(ctx context.Context, userID uuid.UUID) (int, error)

	// Query operations
	GetAll(ctx context.Context, limit, offset int) ([]*models.UserProfile, error)
	Count(ctx context.Context) (int64, error)

	// Validation operations
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	AccountStatusValidation(ctx context.Context, userID uuid.UUID) (bool, error)

	// Helper methods
	VerifyPassword(inputPassword, storedHash string) bool
	HashPassword(password string) (string, error)
	ValidatePasswordStrength(password string) error
}
