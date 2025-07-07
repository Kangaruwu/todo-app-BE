package service

import (
	"context"
	"fmt"

	"go-backend-todo/internal/models"
	"go-backend-todo/internal/repository/todo"

	"github.com/google/uuid"
)

// TodoService interface defines business logic for todo
type TodoService interface {
	CreateTodo(ctx context.Context, req models.CreateTodoRequest, userID uuid.UUID) (*models.Todo, error)
	GetTodoByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Todo, error)
	UpdateTodo(ctx context.Context, id uuid.UUID, req models.UpdateTodoRequest, userID uuid.UUID) (*models.Todo, error)
	DeleteTodo(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetUserTodos(ctx context.Context, userID uuid.UUID, filter models.TodoFilter) ([]*models.Todo, error)
	GetTodosWithPagination(ctx context.Context, userID uuid.UUID, limit, offset int, completed *bool) ([]*models.Todo, int64, error)
	MarkTodosAsCompleted(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error
	DeleteCompletedTodos(ctx context.Context, userID uuid.UUID) error
}

// todoService implementation of TodoService interface
type todoService struct {
	todoRepo todo_repository.TodoRepository
}

// NewTodoService creates a new instance of todo service
func NewTodoService(todoRepo todo_repository.TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

// CreateTodo creates a new todo
func (s *todoService) CreateTodo(ctx context.Context, req models.CreateTodoRequest, userID uuid.UUID) (*models.Todo, error) {
	// Validate input
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	todo := &models.Todo{
		Title:       req.Title,
		Deadline:    req.Deadline,
		Completed:   false,
		UserID:      userID,
	}

	err := s.todoRepo.Create(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

// GetTodoByID retrieves a todo by ID and checks ownership
func (s *todoService) GetTodoByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Todo, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	// Check ownership
	if todo.UserID != userID {
		return nil, fmt.Errorf("todo not found or access denied")
	}

	return todo, nil
}

// UpdateTodo updates a todo
func (s *todoService) UpdateTodo(ctx context.Context, id uuid.UUID, req models.UpdateTodoRequest, userID uuid.UUID) (*models.Todo, error) {
	// Retrieve the current todo and check ownership
	todo, err := s.GetTodoByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided in the request
	if req.Title != nil {
		if *req.Title == "" {
			return nil, fmt.Errorf("title cannot be empty")
		}
		todo.Title = *req.Title
	}

	if req.Deadline != nil {
		todo.Deadline = *req.Deadline
	}

	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	err = s.todoRepo.Update(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return todo, nil
}

// DeleteTodo deletes a todo
func (s *todoService) DeleteTodo(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Check ownership before deleting
	_, err := s.GetTodoByID(ctx, id, userID)
	if err != nil {
		return err
	}

	err = s.todoRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

// GetUserTodos retrieves todos for a user with filter
func (s *todoService) GetUserTodos(ctx context.Context, userID uuid.UUID, filter models.TodoFilter) ([]*models.Todo, error) {
	filter.UserID = userID

	todos, err := s.todoRepo.GetByUserID(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	return todos, nil
}

// GetTodosWithPagination retrieves todos with pagination
func (s *todoService) GetTodosWithPagination(ctx context.Context, userID uuid.UUID, limit, offset int, completed *bool) ([]*models.Todo, int64, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max 100 items per page
	}

	filter := models.TodoFilter{
		UserID:    userID,
		Completed: completed,
		Limit:     limit,
		Offset:    offset,
	}

	todos, err := s.todoRepo.GetByUserID(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get todos: %w", err)
	}

	// Retrieve total count
	total, err := s.todoRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}

	return todos, total, nil
}

// MarkTodosAsCompleted marks multiple todos as completed
func (s *todoService) MarkTodosAsCompleted(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error {
	// TODO: Check ownership of all todos (can be optimized with batch query)
	for _, id := range ids {
		_, err := s.GetTodoByID(ctx, id, userID)
		if err != nil {
			return fmt.Errorf("invalid todo ID %s: %w", id, err)
		}
	}

	err := s.todoRepo.MarkAsCompleted(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to mark todos as completed: %w", err)
	}

	return nil
}

// DeleteCompletedTodos deletes all completed todos for a user
func (s *todoService) DeleteCompletedTodos(ctx context.Context, userID uuid.UUID) error {
	err := s.todoRepo.DeleteCompleted(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete completed todos: %w", err)
	}

	return nil
}
