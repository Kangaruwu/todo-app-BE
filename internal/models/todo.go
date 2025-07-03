package models

import (
	"time"

	"github.com/google/uuid"
)

// Todo model represents a todo item
// It includes fields for ID, title, description, completion status, timestamps, and associated user ID
type Todo struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Deadline    time.Time `json:"deadline" db:"deadline"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
}

// CreateTodoRequest struct represents the request to create a new todo
type CreateTodoRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Deadline    time.Time `json:"deadline" validate:"required"`
}

// UpdateTodoRequest struct represents the request to update an existing todo
type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Deadline    *time.Time `json:"deadline,omitempty" validate:"omitempty"`
	Completed   *bool      `json:"completed,omitempty"`
}

// TodoFilter struct represents the filter for querying todos
type TodoFilter struct {
	UserID    uuid.UUID `json:"user_id"`
	Completed *bool     `json:"completed,omitempty"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}
