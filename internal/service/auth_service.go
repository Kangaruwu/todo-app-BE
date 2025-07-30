package service

import (
	"context"
	"go-backend-todo/internal/config"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/utils"
	"log"
	"time"

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
	// Create timeout context for login operation
	loginCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeouts.AuthTimeout.LoginTimeout)*time.Second)
	defer cancel()

	user, err := s.authRepo.Login(loginCtx, req)
	if err != nil {
		// Check if error is due to timeout
		if loginCtx.Err() == context.DeadlineExceeded {
			log.Println("Login operation timed out")
			return nil, utils.ErrTimeout("Login operation timed out")
		}

		log.Println("Error during login:", err)
		if err.Error() == "invalid credentials" {
			return nil, utils.ErrInvalidCredentials("Invalid email or password")
		}
		return nil, err
	}

	// Create separate timeout context for token version increment
	tokenCtx, tokenCancel := context.WithTimeout(context.Background(), time.Duration(s.config.Timeouts.UserTimeout.TokenVersionTimeout)*time.Second)
	defer tokenCancel()

	err = s.userRepo.IncrementTokenVersion(tokenCtx, user.UserID)
	if err != nil {
		if tokenCtx.Err() == context.DeadlineExceeded {
			log.Printf("Token version increment timed out for user %s", user.UserID)
		} else {
			log.Printf("Failed to increment token version for user %s: %v", user.UserID, err)
		}
		// Don't fail login, just log the error
	}

	// Create timeout context for getting updated user
	userCtx, userCancel := context.WithTimeout(context.Background(), time.Duration(s.config.Timeouts.UserTimeout.GetUserTimeout)*time.Second)
	defer userCancel()

	updatedUser, err := s.userRepo.GetByID(userCtx, user.UserID)
	if err != nil {
		if userCtx.Err() == context.DeadlineExceeded {
			log.Printf("Get updated user timed out for user %s", user.UserID)
		} else {
			log.Printf("Failed to get updated user after token increment: %v", err)
		}
		return user, nil // Return original user if can't get updated one
	}

	return updatedUser, nil
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest, verificationToken string) error {
	// Create timeout context for entire registration process
	registerCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeouts.AuthTimeout.RegisterTimeout)*time.Second)
	defer cancel()

	// Check email existence with timeout
	emailCtx, emailCancel := context.WithTimeout(registerCtx, time.Duration(s.config.Timeouts.UserTimeout.EmailAlreadyExistsTimeout)*time.Second)
	defer emailCancel()

	exists, err := s.userRepo.EmailExists(emailCtx, req.Email)
	if err != nil {
		if emailCtx.Err() == context.DeadlineExceeded {
			log.Printf("Email existence check timed out for email: %s", req.Email)
			return utils.ErrTimeout("Email existence check timed out")
		}
		return err
	}
	if exists {
		return utils.ErrEmailAlreadyExists(req.Email)
	}

	// Check username existence with timeout
	usernameCtx, usernameCancel := context.WithTimeout(registerCtx, time.Duration(s.config.Timeouts.UserTimeout.UsernameExistsTimeout)*time.Second)
	defer usernameCancel()

	exists, err = s.userRepo.UsernameExists(usernameCtx, req.Username)
	if err != nil {
		if usernameCtx.Err() == context.DeadlineExceeded {
			log.Printf("Username existence check timed out for username: %s", req.Username)
			return utils.ErrTimeout("Username existence check timed out")
		}
		return err
	}
	if exists {
		return utils.ErrUsernameAlreadyExists(req.Username)
	}

	// Password strength validation (no timeout needed - local operation)
	if err := s.userRepo.ValidatePasswordStrength(req.Password); err != nil {
		return err
	}

	// Create user account with timeout
	createCtx, createCancel := context.WithTimeout(registerCtx, time.Duration(s.config.Timeouts.UserTimeout.CreateUserTimeout)*time.Second)
	defer createCancel()

	err = s.userRepo.Create(createCtx, req, verificationToken)
	if err != nil {
		if createCtx.Err() == context.DeadlineExceeded {
			log.Printf("User creation timed out for email: %s", req.Email)
			return utils.ErrTimeout("User creation timed out")
		}
		log.Println("Error registering user:", err)
		return err
	}

	// Send confirmation email with timeout (use background context to not inherit parent timeout)
	emailSendCtx, emailSendCancel := context.WithTimeout(context.Background(), time.Duration(s.config.Timeouts.EmailTimeout.EmailSendTimeout)*time.Second)
	defer emailSendCancel()

	err = s.emailService.SendVerificationEmail(emailSendCtx, req.Username, req.Email, verificationToken)
	if err != nil {
		if emailSendCtx.Err() == context.DeadlineExceeded {
			log.Printf("Verification email send timed out for email: %s", req.Email)
			// Don't fail registration if email fails - user is already created
			log.Printf("User %s registered successfully but email failed to send", req.Email)
		} else {
			log.Printf("Error sending verification email to %s: %v", req.Email, err)
		}
		// Don't return error here - user is already registered
	}

	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, verificationToken string) error {
	// Create timeout context for email verification
	verifyCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeouts.EmailTimeout.VerifyEmailTimeout)*time.Second)
	defer cancel()

	createdAt, err := s.authRepo.GetTokenCreationTime(verifyCtx, verificationToken, true)
	if err != nil {
		if verifyCtx.Err() == context.DeadlineExceeded {
			log.Printf("Token validation timed out for token: %s", verificationToken)
			return utils.ErrTimeout("Token validation timed out")
		}
		return err
	}

	if createdAt.Add(time.Duration(s.config.Token.VerifyEmailTokenTTL) * time.Minute).Before(time.Now()) {
		log.Printf("Token is too old: %s", verificationToken)
		return utils.ErrInvalidCredentials("Token has expired")
	}

	// Use the same timeout context for email verification
	err = s.authRepo.VerifyEmail(verifyCtx, verificationToken)
	if err != nil {
		if verifyCtx.Err() == context.DeadlineExceeded {
			log.Printf("Email verification operation timed out for token: %s", verificationToken)
			return utils.ErrTimeout("Email verification timed out")
		}
		return err
	}

	return nil
}

func (s *authService) RecoverPassword(ctx context.Context, req *models.RecoverPasswordRequest, recoverToken string) error {
	// Create timeout context for password recovery
	recoverCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeouts.AuthTimeout.RecoverPasswordTimeout)*time.Second)
	defer cancel()

	// TODO: Implement password recovery logic here
	// For now, we'll just log the operation with timeout context

	// When implemented, send recovery email with timeout
	// emailCtx, emailCancel := context.WithTimeout(recoverCtx, time.Duration(s.config.Timeouts.EmailTimeout.EmailSendTimeout)*time.Second)
	// defer emailCancel()

	// err := s.emailService.SendPasswordResetEmail(emailCtx, req.Email, recoverToken)
	// if err != nil {
	// 	if emailCtx.Err() == context.DeadlineExceeded {
	// 		log.Printf("Recovery email send timed out for email: %s", req.Email)
	// 		return utils.ErrTimeout("Password recovery email send timed out")
	// 	}
	// 	log.Printf("Error sending recovery email to %s: %v", req.Email, err)
	// 	return err
	// }

	// Check if context timed out
	if recoverCtx.Err() == context.DeadlineExceeded {
		log.Printf("Password recovery operation timed out for email: %s", req.Email)
		return utils.ErrTimeout("Password recovery operation timed out")
	}

	log.Printf("Password recovery requested for email: %s", req.Email)
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	// Create timeout context for password reset
	resetCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeouts.AuthTimeout.ResetPasswordTimeout)*time.Second)
	defer cancel()

	createdAt, err := s.authRepo.GetTokenCreationTime(resetCtx, req.Token, false)
	if err != nil {
		if resetCtx.Err() == context.DeadlineExceeded {
			log.Printf("Token validation timed out for reset token: %s", req.Token)
			return utils.ErrTimeout("Token validation timed out")
		}
		return err
	}

	if createdAt.Add(time.Duration(s.config.Token.RecoverPasswordTokenTTL) * time.Minute).Before(time.Now()) {
		log.Printf("Token is too old: %s", req.Token)
		return utils.ErrInvalidCredentials("Token has expired")
	}

	// Use the same timeout context for password reset
	err = s.authRepo.ResetPassword(resetCtx, req)
	if err != nil {
		if resetCtx.Err() == context.DeadlineExceeded {
			log.Printf("Password reset operation timed out for token: %s", req.Token)
			return utils.ErrTimeout("Password reset timed out")
		}
		return err
	}

	return nil
}
