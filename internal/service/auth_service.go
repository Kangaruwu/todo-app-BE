package service

import (
	"context"
	"log"
	"time"
	"go-backend-todo/internal/utils"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/config"

	auth_repository "go-backend-todo/internal/repository/auth"
	user_repository "go-backend-todo/internal/repository/user"
)

type AuthService interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error)
	Register(ctx context.Context, req *models.RegisterRequest, verificationToken string) error

	VerifyEmail(ctx context.Context, verificationToken string) error
	RecoverPassword(ctx context.Context, req *models.RecoverPasswordRequest, recoverToken string) error
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error
}

type authService struct {
	userRepo     user_repository.UserRepository
	authRepo     auth_repository.AuthRepository
	emailService EmailService
	config       *config.Config
}

func NewAuthService(userRepo user_repository.UserRepository, authRepo auth_repository.AuthRepository, emailService EmailService, cfg *config.Config) AuthService {
	return &authService{
		userRepo:     userRepo,
		authRepo:     authRepo,
		emailService: emailService,
		config:       cfg,
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

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest, verificationToken string) error {
	// Condition checks
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return utils.ErrEmailAlreadyExists(req.Email)
	}

	exists, err = s.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		return err
	}
	if exists {
		return utils.ErrUsernameAlreadyExists(req.Username)
	}

	// Create user account
	err = s.userRepo.Create(ctx, req, verificationToken)
	if err != nil {
		log.Println("Error registering user:", err)
		return err
	}

	// Send confirmation email
	err = s.emailService.SendVerificationEmail(ctx, req.Username, req.Email, verificationToken)
	if err != nil {
		log.Println("Error sending verification email:", err)
		return err
	}

	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, verificationToken string) error {
    createdAt, err := s.authRepo.GetTokenCreationTime(ctx, verificationToken, true)
    if err != nil {
        return err
    }
    
    if createdAt.Add(time.Duration(s.config.Token.VerifyEmailTokenTTL) * time.Minute).Before(time.Now()) {
        log.Printf("Token is too old: %s", verificationToken)
        return utils.ErrInvalidCredentials("Token has expired")
    }
	return s.authRepo.VerifyEmail(ctx, verificationToken)
}

func (s *authService) RecoverPassword(ctx context.Context, req *models.RecoverPasswordRequest, recoverToken string) error {
	// TODO: Send recovery email with the token
	// if err := s.emailService.SendPasswordResetEmail(ctx, req.Email, ..., recoverToken); err != nil {
	// 	log.Println("Error sending recovery email:", err)
	// 	return err
	// }
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
    createdAt, err := s.authRepo.GetTokenCreationTime(ctx, req.Token, false)
    if err != nil {
        return err
    }
    
    if createdAt.Add(time.Duration(s.config.Token.RecoverPasswordTokenTTL) * time.Minute).Before(time.Now()) {
        log.Printf("Token is too old: %s", req.Token)
        return utils.ErrInvalidCredentials("Token has expired")
    }
	
	return s.authRepo.ResetPassword(ctx, req)
}

