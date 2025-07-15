package auth_repository

import (
	"context"
	"time"
	"go-backend-todo/internal/models"
)

type AuthRepository interface {
	ValidateCredentials(ctx context.Context, email, password string) (*models.UserAccount, error)
	VerifyEmail(ctx context.Context, token string) error
	RecoverPassword(ctx context.Context, email string) (*models.UserAccount, error)
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error
	Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error)
	GetTokenCreationTime(ctx context.Context, token string, isVerifyToken bool) (time.Time, error)
}
