package user_repository

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
func (u *userRepository) Create(ctx context.Context, user *models.UserAccount) error {
	return nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserAccount, error) {
	return nil, nil
}
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*models.UserAccount, error) {
	return nil, nil
}
func (u *userRepository) GetByUsername(ctx context.Context, username string) (*models.UserAccount, error) {
	return nil, nil
}
func (u *userRepository) Update(ctx context.Context, user *models.UserAccount) error {
	return nil
}
func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// Query operations
func (u *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.UserAccount, error) {
	query := "SELECT * FROM user_account;"
	pool := db.GetPool()
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Println("BRUHHHHHHHHHH:", err)
	}
	defer rows.Close()
	var users []*models.UserAccount
	for rows.Next() {
		var user models.UserAccount
		err := rows.Scan(
			&user.UserID, &user.UserRole,
		)	
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
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
