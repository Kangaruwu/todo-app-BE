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
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	Register(ctx context.Context, user *models.RegisterRequest) error
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

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	loginResponse, err := s.authRepo.Login(ctx, req)
	if err != nil {
		log.Println("Error during login:", err)
		if err.Error() == "invalid credentials" {
			return nil, utils.ErrInvalidCredentials("Invalid email or password")
		}
		return nil, err
	}
	return loginResponse, nil
}

func (s *authService) Register(ctx context.Context, user *models.RegisterRequest) error {
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Println("Error registering user:", err)
		return err
	}
	return nil
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