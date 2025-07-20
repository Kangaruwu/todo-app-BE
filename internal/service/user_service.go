package service

import (
	"context"
	"go-backend-todo/internal/models"
	user_repository "go-backend-todo/internal/repository/user"
	"go-backend-todo/internal/utils"
	"log"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int64, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req models.UpdateProfileRequest) (*models.UserProfile, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetUserStats(ctx context.Context) (*models.UserStatsResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error
}

type userService struct {
	userRepo user_repository.UserRepository
}

func NewUserService(userRepo user_repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int64, error) {
	// TODO: Implement get all users with pagination
	return nil, 0, nil
}

func (s *userService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req models.UpdateProfileRequest) (*models.UserProfile, error) {
	// TODO: Implement update user profile
	return nil, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// TODO: Implement delete user
	return nil
}

func (s *userService) GetUserStats(ctx context.Context) (*models.UserStatsResponse, error) {
	// TODO: Implement user statistics
	stats := &models.UserStatsResponse{
		TotalUsers:         0,
		ActiveUsers:        0,
		PendingUsers:       0,
		RegisteredToday:    0,
		RegisteredThisWeek: 0,
	}
	return stats, nil
}

// ChangePassword changes the user's password
func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error {
    log.Println("Changing password for user:", userID, "Request:", req)
    // 1. Validate input structure
    if req.NewPassword != req.ConfirmPassword {
        return utils.ErrInvalidInput("New password and confirm password do not match")
    }

    // 2. Validate password strength
    if err := s.userRepo.ValidatePasswordStrength(req.NewPassword); err != nil {
        return err
    }

    // 3. Verify current password
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return utils.ErrInternalServerError("Failed to get user")
    }

    if !s.userRepo.VerifyPassword(req.CurrentPassword, user.PasswordHash) {
        return utils.ErrInvalidCredentials("Current password is incorrect")
    }

    // // 4. Check password history (prevent reuse)
    // if err := s.checkPasswordHistory(ctx, userID, req.NewPassword); err != nil {
    //     return err
    // }

    // 5. Hash and update password
    hashedPassword, err := s.userRepo.HashPassword(req.NewPassword)
    if err != nil {
        return utils.ErrInternalServerError("Failed to hash password")
    }

    return s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
}

