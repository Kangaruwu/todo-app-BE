package auth_repository

import (
	"context"
	"go-backend-todo/internal/models"
	"log"
	"net/url"
	"time"

	"go-backend-todo/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	// Check if context is already cancelled/timed out
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	query := "SELECT user_id, user_name, password_hash, email_validation_status FROM user_account WHERE email_address = $1;"

	var userID uuid.UUID
	var userName string
	var passwordHash string
	var emailValidationStatus string

	err := a.db.QueryRow(ctx, query, email).Scan(&userID, &userName, &passwordHash, &emailValidationStatus)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("ValidateCredentials operation timed out")
			return nil, utils.ErrTimeout("Credential validation timed out")
		}
		if ctx.Err() == context.Canceled {
			log.Println("ValidateCredentials operation was cancelled")
			return nil, utils.ErrInternalServerError("Credential validation was cancelled")
		}
		log.Println("Error validating credentials:", err)
		return nil, utils.ErrInvalidCredentials("Invalid credentials")
	}

	// Password comparison (local operation, no need for context timeout check)
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		log.Println("Invalid password:", err)
		return nil, utils.ErrInvalidCredentials("Invalid credentials")
	}

	return &models.UserAccount{
		UserID:   userID,
		UserRole: models.UserRoleEnum("user"), // Default role
	}, nil
}

func (a *authRepository) VerifyEmail(ctx context.Context, token string) error {
	// Check if context is already cancelled/timed out
	if ctx.Err() != nil {
		return ctx.Err()
	}

	query := "UPDATE user_account SET email_validation_status = 'confirmed'::email_validation_status_enum WHERE verification_token = $1;"
	result, err := a.db.Exec(ctx, query, token)
	if err != nil {
		// Check if error is due to context timeout/cancellation
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("VerifyEmail operation timed out")
			return utils.ErrTimeout("Email verification timed out")
		}
		if ctx.Err() == context.Canceled {
			log.Println("VerifyEmail operation was cancelled")
			return utils.ErrInternalServerError("Email verification was cancelled")
		}

		log.Println("Error verifying email:", err)
		return utils.ErrInternalServerError("Failed to verify email")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrInvalidCredentials("Invalid verification token")
	}

	return nil
}

func (a *authRepository) GetTokenCreationTime(ctx context.Context, token string, isVerifyToken bool) (time.Time, error) {
	// Check if context is already cancelled/timed out
	if ctx.Err() != nil {
		return time.Time{}, ctx.Err()
	}

	decodedToken, err := url.QueryUnescape(token)
	if err != nil {
		return time.Time{}, utils.ErrInvalidCredentials("Invalid verification token format")
	}

	var query string
	if isVerifyToken {
		query = "SELECT verification_token_generation_time FROM user_account WHERE verification_token = $1;"
	} else {
		query = "SELECT password_recovery_token_generation_time FROM user_account WHERE password_recovery_token = $1;"
	}

	var createdAt time.Time
	err = a.db.QueryRow(ctx, query, decodedToken).Scan(&createdAt)
	if err != nil {
		// Check if error is due to context timeout/cancellation
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("GetTokenCreationTime operation timed out")
			return time.Time{}, utils.ErrTimeout("Token validation timed out")
		}
		if ctx.Err() == context.Canceled {
			log.Println("GetTokenCreationTime operation was cancelled")
			return time.Time{}, utils.ErrInternalServerError("Token validation was cancelled")
		}

		log.Println("Error getting token creation time:", err)
		return time.Time{}, utils.ErrInvalidCredentials("Invalid token")
	}

	return createdAt, nil
}

func (a *authRepository) RecoverPassword(ctx context.Context, email string) (*models.UserAccount, error) {
	// Check if context is already cancelled/timed out
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	query := "SELECT user_id, user_role FROM user_account WHERE email_address = $1;"

	var userID uuid.UUID
	var userRole string

	err := a.db.QueryRow(ctx, query, email).Scan(&userID, &userRole)
	if err != nil {
		// Check if error is due to context timeout/cancellation
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("RecoverPassword lookup timed out")
			return nil, utils.ErrTimeout("Password recovery lookup timed out")
		}
		if ctx.Err() == context.Canceled {
			log.Println("RecoverPassword lookup was cancelled")
			return nil, utils.ErrInternalServerError("Password recovery lookup was cancelled")
		}

		log.Println("Error during password recovery lookup:", err)
		return nil, utils.ErrUserNotFound("Email not found")
	}

	return &models.UserAccount{
		UserID:   userID,
		UserRole: models.UserRoleEnum(userRole),
	}, nil
}

func (a *authRepository) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	// Check if context is already cancelled/timed out
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Hash the new password first (local operation)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return utils.ErrInternalServerError("Failed to process password")
	}

	query := "UPDATE user_account SET password_hash = $1 WHERE password_recovery_token = $2;"
	result, err := a.db.Exec(ctx, query, string(hashedPassword), req.Token)
	if err != nil {
		// Check if error is due to context timeout/cancellation
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("ResetPassword operation timed out")
			return utils.ErrTimeout("Password reset timed out")
		}
		if ctx.Err() == context.Canceled {
			log.Println("ResetPassword operation was cancelled")
			return utils.ErrInternalServerError("Password reset was cancelled")
		}

		log.Println("Error resetting password:", err)
		return utils.ErrInternalServerError("Failed to reset password")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return utils.ErrInvalidCredentials("Invalid recovery token or token has expired")
	}

	return nil
}

func (a *authRepository) Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error) {
	// Check if context is already cancelled/timed out
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	query := "SELECT user_id, user_name, user_role, password_hash, email_validation_status FROM user_account WHERE email_address = $1;"

	var passwordHash string
	var userID uuid.UUID
	var userName string
	var userRole string
	var emailValidationStatus string

	err := a.db.QueryRow(ctx, query, req.Email).Scan(&userID, &userName, &userRole, &passwordHash, &emailValidationStatus)
	if err != nil {
		// Check if error is due to context timeout/cancellation
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Login database query timed out")
			return nil, utils.ErrTimeout("Login operation timed out")
		}
		if ctx.Err() == context.Canceled {
			log.Println("Login database query was cancelled")
			return nil, utils.ErrInternalServerError("Login operation was cancelled")
		}

		log.Println("Error during login:", err)
		return nil, utils.ErrInvalidCredentials("Invalid email or password")
	}

	// Password comparison (local operation, no need for context timeout check)
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		log.Println("Invalid credentials:", err)
		return nil, utils.ErrInvalidCredentials("Invalid email or password")
	}

	return &models.UserProfile{
		UserID:   userID,
		Username: userName,
		Email:    req.Email,
		Role:     models.UserRoleEnum(userRole),
		Status:   models.EmailValidationStatusEnum(emailValidationStatus),
	}, nil
}
