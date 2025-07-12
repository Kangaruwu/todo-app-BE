package models

import (
	"time"

	"github.com/google/uuid"
)

// Todo model represents a todo item
// It includes fields for ID, title, description, completion status, timestamps, and associated user ID
type Todo struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Completed bool      `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Deadline  time.Time `json:"deadline" db:"deadline"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
}

// CreateTodoRequest struct represents the request to create a new todo
type CreateTodoRequest struct {
	Title    string    `json:"title" validate:"required,min=1,max=255"`
	Deadline time.Time `json:"deadline" validate:"required"`
}

// UpdateTodoRequest struct represents the request to update an existing todo
type UpdateTodoRequest struct {
	Title     *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Deadline  *time.Time `json:"deadline,omitempty" validate:"omitempty"`
	Completed *bool      `json:"completed,omitempty"`
}

// TodoFilter struct represents the filter for querying todos
type TodoFilter struct {
	UserID    uuid.UUID `json:"user_id"`
	Completed *bool     `json:"completed,omitempty"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

// TodoListResponse represents paginated todo list response
type TodoListResponse struct {
	Todos  []Todo `json:"todos"`
	Total  int64  `json:"total"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// TodoStatsResponse represents todo statistics
type TodoStatsResponse struct {
	TotalTodos     int64 `json:"total_todos"`
	CompletedTodos int64 `json:"completed_todos"`
	PendingTodos   int64 `json:"pending_todos"`
	OverdueTodos   int64 `json:"overdue_todos"`
	TodayTodos     int64 `json:"today_todos"`
	ThisWeekTodos  int64 `json:"this_week_todos"`
}

// Priority enum for todo priority levels
type TodoPriority string

const (
	LowPriority    TodoPriority = "low"
	MediumPriority TodoPriority = "medium"
	HighPriority   TodoPriority = "high"
	UrgentPriority TodoPriority = "urgent"
)

// Category represents todo categories
type TodoCategory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Color     string    `json:"color" db:"color"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Enhanced Todo model with priority and category
type TodoEnhanced struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	Title       string        `json:"title" db:"title"`
	Description string        `json:"description" db:"description"`
	Completed   bool          `json:"completed" db:"completed"`
	Priority    TodoPriority  `json:"priority" db:"priority"`
	CategoryID  *uuid.UUID    `json:"category_id" db:"category_id"`
	Category    *TodoCategory `json:"category,omitempty"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	Deadline    *time.Time    `json:"deadline" db:"deadline"`
	UserID      uuid.UUID     `json:"user_id" db:"user_id"`
}

// CreateTodoEnhancedRequest for enhanced todo creation
type CreateTodoEnhancedRequest struct {
	Title       string       `json:"title" validate:"required,min=1,max=255" example:"Complete project proposal"`
	Description string       `json:"description" validate:"max=1000" example:"Write and submit the quarterly project proposal"`
	Priority    TodoPriority `json:"priority" validate:"required,oneof=low medium high urgent" example:"high"`
	CategoryID  *uuid.UUID   `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Deadline    *time.Time   `json:"deadline" example:"2025-07-15T10:00:00Z"`
}

// UpdateTodoEnhancedRequest for enhanced todo updates
type UpdateTodoEnhancedRequest struct {
	Title       *string       `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string       `json:"description,omitempty" validate:"omitempty,max=1000"`
	Priority    *TodoPriority `json:"priority,omitempty" validate:"omitempty,oneof=low medium high urgent"`
	CategoryID  *uuid.UUID    `json:"category_id,omitempty"`
	Deadline    *time.Time    `json:"deadline,omitempty"`
	Completed   *bool         `json:"completed,omitempty"`
}

// CreateCategoryRequest for category creation
type CreateCategoryRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=100" example:"Work"`
	Color string `json:"color" validate:"required,hexcolor" example:"#FF5733"`
}
