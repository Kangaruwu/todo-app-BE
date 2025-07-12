package user_repository

import (
	"context"
	"go-backend-todo/internal/db"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/utils"
	"go-backend-todo/internal/repository/auth"
	"time"

	"log"
    "golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// CRUD operations
func (u *userRepository) Create(ctx context.Context, req *models.RegisterRequest) error {
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
	err = auth_repository.SendVerificationEmail(ctx, req.Username, req.Email, confirmationToken)
	if err != nil {
		log.Println("Error sending verification email:", err)
		return err
	}

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
	return false, nil
}
func (u *userRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return false, nil
}
