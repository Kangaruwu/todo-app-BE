package user_repository

import (
	"context"

	"go-backend-todo/internal/models"

	"github.com/google/uuid"
)

// UserRepository interface defines methods for interacting with user data
type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *models.UserAccount) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.UserAccount, error)
	GetByEmail(ctx context.Context, email string) (*models.UserAccount, error)
	GetByUsername(ctx context.Context, username string) (*models.UserAccount, error)
	Update(ctx context.Context, user *models.UserAccount) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetAll(ctx context.Context, limit, offset int) ([]*models.UserAccount, error)
	Count(ctx context.Context) (int64, error)

	// Validation operations
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
}
