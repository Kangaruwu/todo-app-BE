package user_repository

import (
	"context"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/utils"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// CRUD operations
func (u *userRepository) Create(ctx context.Context, req *models.RegisterRequest, verificationToken string) error {
	pw_hash, err := u.HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	query := `
		INSERT INTO user_account (
			user_name, 
			user_role,
			password_hash, 
			email_address, 
			verification_token, 
			verification_token_generation_time, 
			email_validation_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err = u.db.Exec(ctx, query, req.Username, "user", pw_hash, req.Email, verificationToken, time.Now(), "pending")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	query := "SELECT user_id, user_name, password_hash, email_address, user_role, email_validation_status, COALESCE(token_version, 1), created_at, updated_at FROM user_account WHERE user_id = $1;"
	row := u.db.QueryRow(ctx, query, id)
	var user models.UserProfile
	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Role,
		&user.Status,
		&user.TokenVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, utils.ErrUserNotFound("User not found")
		}
		log.Println("Error fetching user by ID:", err)
		return nil, err
	}
	return &user, nil
}
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*models.UserProfile, error) {
	return nil, nil
}
func (u *userRepository) GetByUsername(ctx context.Context, username string) (*models.UserProfile, error) {
	return nil, nil
}
func (u *userRepository) Update(ctx context.Context, user *models.UserAccount) error {
	return nil
}
func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// Query operations
func (u *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.UserProfile, error) {
	return nil, nil
}
func (u *userRepository) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

// Validation operations
func (u *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_account WHERE email_address = $1);"
	var exists bool
	err := u.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Println("Error checking email existence:", err)
		return false, err
	}
	return exists, nil
}

func (u *userRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_account WHERE user_name = $1);"
	var exists bool
	err := u.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		log.Println("Error checking username existence:", err)
		return false, err
	}
	return exists, nil
}

func (u *userRepository) AccountStatusValidation(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := "SELECT email_validation_status FROM user_account WHERE user_id = $1;"
	var status string
	err := u.db.QueryRow(ctx, query, userID).Scan(&status)
	if err != nil {
		log.Println("Error checking account status:", err)
		return false, err
	}

	if status == "active" {
		return true, nil
	}
	return false, utils.ErrAccountNotActive("User account is not active")
}

func (u *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return utils.ErrInternalServerError("Failed to start transaction")
	}
	defer tx.Rollback(ctx)

	// Update password and increment token version in same transaction
	query := `UPDATE user_account 
			  SET password_hash = $1, 
				  token_version = token_version + 1,
				  updated_at = CURRENT_TIMESTAMP 
			  WHERE user_id = $2`

	result, err := tx.Exec(ctx, query, newPassword, userID)
	if err != nil {
		log.Println("Error updating password:", err)
		return utils.ErrInternalServerError("Failed to update password")
	}

	if result.RowsAffected() == 0 {
		return utils.ErrUserNotFound("User not found")
	}

	return tx.Commit(ctx)
}

// Token version operations
func (u *userRepository) IncrementTokenVersion(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE user_account SET token_version = token_version + 1 WHERE user_id = $1`
	result, err := u.db.Exec(ctx, query, userID)
	if err != nil {
		log.Println("Error incrementing token version:", err)
		return utils.ErrInternalServerError("Failed to increment token version")
	}

	if result.RowsAffected() == 0 {
		return utils.ErrUserNotFound("User not found")
	}

	return nil
}

func (u *userRepository) GetTokenVersion(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(token_version, 1) FROM user_account WHERE user_id = $1`
	var version int
	err := u.db.QueryRow(ctx, query, userID).Scan(&version)
	if err != nil {
		log.Println("Error getting token version:", err)
		return 0, utils.ErrInternalServerError("Failed to get token version")
	}
	return version, nil
}

func (u *userRepository) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return utils.ErrInvalidInput("Password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return utils.ErrInvalidInput("Password must not exceed 128 characters")
	}

	// Check for required character types
	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return utils.ErrInvalidInput("Password must contain at least one lowercase letter, one uppercase letter, one digit, and one special character")
	}

	return nil
}

func (u *userRepository) VerifyPassword(inputPassword, storedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPassword))
	return err == nil
}

func (u *userRepository) HashPassword(password string) (string, error) {
	salt := utils.RandInRange(bcrypt.MinCost, bcrypt.MaxCost)
	pw_hash, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	if err != nil {
		return "", utils.ErrInternalServerError("Failed to hash password")
	}
	return string(pw_hash), nil
}
