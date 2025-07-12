package user_repository

import (
	"context"
	"go-backend-todo/internal/db"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/repository/auth"
	"go-backend-todo/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// CRUD operations
func (u *userRepository) Create(ctx context.Context, req *models.RegisterRequest) error {
	exists, err := u.EmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return utils.ErrEmailAlreadyExists(req.Email)
	}

	exists, err = u.UsernameExists(ctx, req.Username)
	if err != nil {
		return err
	}
	if exists {
		return utils.ErrUsernameAlreadyExists(req.Username)
	}

	salt := utils.RandInRange(bcrypt.MinCost, bcrypt.MaxCost)
	pw_hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), salt)
	if err != nil {
		log.Println(err)
	}

	pool := db.GetPool()
	// Create email validation token
	confirmationToken, err := auth_repository.GenerateEmailValidationToken()
	if err != nil {
		log.Println("Error generating email validation token:", err)
		return err
	}

	query := `
		INSERT INTO user_account (
			user_name, 
			user_role,
			password_hash, 
			hash_algorithm, 
			email_address, 
			confirmation_token, 
			confirmation_token_generation_time, 
			email_validation_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err = pool.Exec(ctx, query, req.Username, "user", pw_hash, "bcrypt", req.Email, confirmationToken, time.Now(), "pending")
	if err != nil {
		log.Println(err)
		return err
	}

	//Send verification email
	// err = auth_repository.SendVerificationEmail(ctx, req.Username, req.Email, confirmationToken)
	// if err != nil {
	// 	log.Println("Error sending verification email:", err)
	// 	return err
	// }

	return nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserAccount, error) {
	return nil, nil
}
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*models.UserAccount, error) {
	return nil, nil
}
func (u *userRepository) GetByUsername(ctx context.Context, username string) (*models.UserAccount, error) {
	return nil, nil
}
func (u *userRepository) Update(ctx context.Context, user *models.UserAccount) error {
	return nil
}
func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// Query operations
func (u *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.UserAccount, error) {
	query := "SELECT * FROM user_account;"
	pool := db.GetPool()
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Println("BRUHHHHHHHHHH:", err)
	}
	defer rows.Close()
	var users []*models.UserAccount
	for rows.Next() {
		var user models.UserAccount
		err := rows.Scan(
			&user.UserID, &user.UserRole,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
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