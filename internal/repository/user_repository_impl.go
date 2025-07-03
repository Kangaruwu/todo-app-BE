package repository

import (
	"context"
	"go-backend-todo/internal/db"
	"go-backend-todo/internal/models"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// CRUD operations
func (u *userRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (email, username, password) VALUES ($1, $2, $3)"

	pool := db.GetPool()
	res, err := pool.Exec(ctx, query,
		user.Email, user.Username, "lol",
	)
	if res.RowsAffected() == 0 {
		log.Println("Cannot create user:", err)
		return err
	}
	log.Printf("Create new user %s successfully!", user.Username)
	return nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return nil, nil
}
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (u *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}
func (u *userRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}
func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// Query operations
func (u *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	return nil, nil
}
func (u *userRepository) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

// Validation operations
func (u *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (u *userRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return false, nil
}
