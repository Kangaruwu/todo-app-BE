package auth_repository

import (
	"context"
	"go-backend-todo/internal/models"
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-backend-todo/internal/utils"
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
	return nil, nil
}

func (a *authRepository) VerifyEmail(ctx context.Context, token string) error {    
    query := "UPDATE user_account SET email_validation_status = 'confirmed'::email_validation_status_enum WHERE verification_token = $1;"
    result, err := a.db.Exec(ctx, query, token)
    if err != nil {
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
        return time.Time{}, utils.ErrInternalServerError("Failed to get token creation time")
    }
    
    return createdAt, nil
}

func (a *authRepository) RecoverPassword(ctx context.Context, email string) (*models.UserAccount, error) {
	return nil, nil
}

func (a *authRepository) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	query := "UPDATE user_account SET password_hash = $1 WHERE recovery_token = $2;"
	_, err := a.db.Exec(ctx, query, req.NewPassword, req.Token)
	if err != nil {
		log.Println("Error resetting password:", err)
		return utils.ErrInternalServerError("Failed to reset password")
	}
	return nil
}

func (a *authRepository) Login(ctx context.Context, req *models.LoginRequest) (*models.UserProfile, error) {
	query := "SELECT user_id, user_name, user_role, password_hash, email_validation_status FROM user_account WHERE email_address = $1;"

	var passwordHash string
	var userID uuid.UUID
	var userName string
	var userRole string
	var emailValidationStatus string
	err := a.db.QueryRow(ctx, query, req.Email).Scan(&userID, &userName, &userRole, &passwordHash, &emailValidationStatus)
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
		UserID:   userID,
		Username: userName,
		Email:    req.Email,
		Role:     models.UserRoleEnum(userRole),
		Status:   models.EmailValidationStatusEnum(emailValidationStatus),
	}, nil
}
