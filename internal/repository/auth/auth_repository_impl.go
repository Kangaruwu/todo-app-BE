package auth_repository

import (
	"context"
	"go-backend-todo/internal/models"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (a *authRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_login_data WHERE email_address = $1);"
	var exists bool
	err := a.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Println("Error checking email existence:", err)
		return false, err
	}
	return exists, nil
}

func (a *authRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_login_data WHERE user_name = $1);"
	var exists bool
	err := a.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		log.Println("Error checking username existence:", err)
		return false, err
	}
	return exists, nil
}

func (a *authRepository) ValidateCredentials(ctx context.Context, email, password string) (*models.UserAccount, error) {
	return nil, nil 
}

func (a *authRepository) CreateUser(ctx context.Context, user *models.CreateUserRequest) (*models.UserAccount, error) {
	return nil, nil
}

func (a *authRepository) ConfirmEmail(ctx context.Context, token string) (*models.UserAccount, error) {
	return nil, nil
}	

func (a *authRepository) RecoverPassword(ctx context.Context, email string) (*models.UserAccount, error) {
	return nil, nil
}

func (a *authRepository) ResetPassword(ctx context.Context, token, newPassword string) (*models.UserAccount, error) {
	return nil, nil
}

