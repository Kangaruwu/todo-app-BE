package auth_repository

import (
	"context"
	"go-backend-todo/internal/models"
)

type AuthRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	ValidateCredentials(ctx context.Context, email, password string) (*models.UserAccount, error)
	CreateUser(ctx context.Context, user *models.CreateUserRequest) (*models.UserAccount, error)
	ConfirmEmail(ctx context.Context, token string) (*models.UserAccount, error)
	RecoverPassword(ctx context.Context, email string) (*models.UserAccount, error)
	ResetPassword(ctx context.Context, token, newPassword string) (*models.UserAccount, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error)
}
