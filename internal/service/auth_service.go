package service

import (
	"context"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/repository/user"
	"go-backend-todo/internal/repository/auth"	
	"go-backend-todo/internal/utils"
	"log"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*models.UserAccount, error)
	Register(ctx context.Context, user *models.UserAccount) (*models.UserAccount, error)
	ValidateEmail(ctx context.Context, email string) (bool, error)
	GenerateEmailValidationToken(ctx context.Context, email string) (string, error)
	ValidateEmailToken(ctx context.Context, token string) (bool, error)
}

type authService struct {
	userRepo user_repository.UserRepository
	authRepo auth_repository.AuthRepository
}
func NewAuthService(userRepo user_repository.UserRepository, authRepo auth_repository.AuthRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.UserAccount, error) {
	return s.authRepo.ValidateCredentials(ctx, email, password)
}

func (s *authService) Register(ctx context.Context, user *models.UserAccount) (*models.UserAccount, error) {
	// return s.authRepo.CreateUser(ctx, user)
	return nil, utils.ErrNotImplemented("Register method not implemented")
}

func (s *authService) ValidateEmail(ctx context.Context, email string) (bool, error) {
	return s.authRepo.EmailExists(ctx, email)
}
func (s *authService) GenerateEmailValidationToken(ctx context.Context, email string) (string, error) {
	return "", utils.ErrNotImplemented("GenerateEmailValidationToken method not implemented")
}

func (s *authService) ValidateEmailToken(ctx context.Context, token string) (bool, error) {
	user, err := s.authRepo.ConfirmEmail(ctx, token)
	if err != nil {
		log.Println("Error validating email token:", err)
		return false, err
	}
	if user == nil {
		return false, nil // Token is invalid or expired
	}
	return true, nil // Token is valid
}