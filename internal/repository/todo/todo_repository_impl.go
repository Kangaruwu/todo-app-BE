package todo_repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-backend-todo/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// todoRepository implementation of TodoRepository interface
type todoRepository struct {
	db *pgxpool.Pool
}

// NewTodoRepository create a new instance of todo repository
func NewTodoRepository(db *pgxpool.Pool) TodoRepository {
	return &todoRepository{db: db}
}

// Create a new todo
func (r *todoRepository) Create(ctx context.Context, todo *models.Todo) error {
	query := `
		INSERT INTO todos (id, title, deadline, completed, created_at, updated_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	todo.ID = uuid.New()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		todo.ID, todo.Title, todo.Deadline, todo.Completed,
		todo.CreatedAt, todo.UpdatedAt, todo.UserID,
	)

	return err
}

// GetByID retrieves a todo by its ID
func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
	query := `
		SELECT id, title, deadline, completed, created_at, updated_at, user_id
		FROM todos
		WHERE id = $1
	`

	var todo models.Todo
	err := r.db.QueryRow(ctx, query, id).Scan(
		&todo.ID, &todo.Title, &todo.Deadline, &todo.Completed,
		&todo.CreatedAt, &todo.UpdatedAt, &todo.UserID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, err
	}

	return &todo, nil
}

// Update a todo
func (r *todoRepository) Update(ctx context.Context, todo *models.Todo) error {
	query := `
		UPDATE todos
		SET title = $2, deadline = $3, completed = $4, updated_at = $5
		WHERE id = $1
	`

	todo.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query,
		todo.ID, todo.Title, todo.Deadline, todo.Completed, todo.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

// Delete a todo
func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

// GetByUserID retrieves todos by user ID with filter
func (r *todoRepository) GetByUserID(ctx context.Context, filter models.TodoFilter) ([]*models.Todo, error) {
	query := `
		SELECT id, title, deadline, completed, created_at, updated_at, user_id
		FROM todos
		WHERE user_id = $1
	`
	args := []interface{}{filter.UserID}
	argIndex := 2

	// Filter by completion status if provided
	if filter.Completed != nil {
		query += fmt.Sprintf(" AND completed = $%d", argIndex)
		args = append(args, *filter.Completed)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	// Add limit and offset
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Deadline, &todo.Completed,
			&todo.CreatedAt, &todo.UpdatedAt, &todo.UserID,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	
	return todos, nil
}

// GetAll retrieves all todos with optional filters
func (r *todoRepository) GetAll(ctx context.Context, filter models.TodoFilter) ([]*models.Todo, error) {
	query := `
		SELECT id, title, deadline, completed, created_at, updated_at, user_id
		FROM todos
		WHERE 1=1 
	`
	args := []interface{}{}
	argIndex := 1

	// Add filter by completion status if provided
	if filter.Completed != nil {
		query += fmt.Sprintf(" AND completed = $%d", argIndex)
		args = append(args, *filter.Completed)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	// Add limit and offset
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Deadline, &todo.Completed,
			&todo.CreatedAt, &todo.UpdatedAt, &todo.UserID,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

// Count counts the number of todos for a user with optional filters
func (r *todoRepository) Count(ctx context.Context, filter models.TodoFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM todos WHERE user_id = $1`
	args := []interface{}{filter.UserID}
	argIndex := 2

	if filter.Completed != nil {
		query += fmt.Sprintf(" AND completed = $%d", argIndex)
		args = append(args, *filter.Completed)
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// MarkAsCompleted marks multiple todos as completed
func (r *todoRepository) MarkAsCompleted(ctx context.Context, ids []uuid.UUID) error {
	query := `
		UPDATE todos
		SET completed = true, updated_at = $1
		WHERE id = ANY($2)
	`

	_, err := r.db.Exec(ctx, query, time.Now(), ids)
	return err
}

// DeleteCompleted deletes all completed todos for a user
func (r *todoRepository) DeleteCompleted(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM todos WHERE user_id = $1 AND completed = true`

	_, err := r.db.Exec(ctx, query, userID)
	return err
}
