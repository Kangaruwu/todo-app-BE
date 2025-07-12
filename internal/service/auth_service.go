package service

import (
	"context"
	"log"
	"go-backend-todo/internal/utils"
	"go-backend-todo/internal/models"
	auth_repository "go-backend-todo/internal/repository/auth"
	user_repository "go-backend-todo/internal/repository/user"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error)
	Register(ctx context.Context, req *models.RegisterRequest) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	ConfirmEmail(ctx context.Context, token string) error
	RecoverPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
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

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error) {
	user, err := s.authRepo.Login(ctx, req)
	if err != nil {
		log.Println("Error during login:", err)
		if err.Error() == "invalid credentials" {
			return nil, utils.ErrInvalidCredentials("Invalid email or password")
		}
		return nil, err
	}
	return user, nil
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) error {
	err := s.userRepo.Create(ctx, req)
	if err != nil {
		log.Println("Error registering user:", err)
		return err
	}
	return nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	// TODO: Implement get user by ID
	return nil, nil
}

func (s *authService) ConfirmEmail(ctx context.Context, token string) error {
	// TODO: Implement email confirmation
	return nil
}

func (s *authService) RecoverPassword(ctx context.Context, email string) error {
	// TODO: Implement password recovery
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// TODO: Implement password reset
	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	// TODO: Implement password change
	return nil
}
