package service

import (
	"context"
	"go-backend-todo/internal/models"
	user_repository "go-backend-todo/internal/repository/user"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int64, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req models.UpdateProfileRequest) (*models.UserProfile, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetUserStats(ctx context.Context) (*models.UserStatsResponse, error)
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
	// TODO: Implement get user by ID
	return nil, nil
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
