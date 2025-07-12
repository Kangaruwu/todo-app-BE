package auth_repository

import (
	"context"
	"go-backend-todo/internal/models"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"go-backend-todo/internal/utils"
	"github.com/google/uuid"
)

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{
		db: db,
	}
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

func (a *authRepository) Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error) {
	query := "SELECT user_id, user_name, user_role, password_hash FROM user_account WHERE email_address = $1;"

	var passwordHash string
	var userID uuid.UUID
	var userName string
	var userRole string
	err := a.db.QueryRow(ctx, query, req.Email).Scan(&userID, &userName, &userRole, &passwordHash)
	if err != nil {
		log.Println("Error during login:", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		log.Println("Invalid credentials:", err)
		return nil, utils.ErrInvalidCredentials("Invalid email or password")
	}

	return &models.UserProfile{
		UserID:    userID,
		Username:  userName,
		Email: req.Email,
		Role:  models.UserRoleEnum(userRole),

	}, nil
}