package repository

import (
	"context"

	"go-backend-todo/internal/models"

	"github.com/google/uuid"
)

// TodoRepository interface defines methods for interacting with todo data
type TodoRepository interface {
	// CRUD operations
	Create(ctx context.Context, todo *models.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByUserID(ctx context.Context, filter models.TodoFilter) ([]*models.Todo, error)
	GetAll(ctx context.Context, filter models.TodoFilter) ([]*models.Todo, error)
	Count(ctx context.Context, filter models.TodoFilter) (int64, error)

	// Bulk operations
	MarkAsCompleted(ctx context.Context, ids []uuid.UUID) error
	DeleteCompleted(ctx context.Context, userID uuid.UUID) error
}
