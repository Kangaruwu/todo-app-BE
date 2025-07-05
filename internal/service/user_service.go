package service

import (
	"github.com/gofiber/fiber/v2"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/repository"
)

type UserService interface {
	CreateUser()
	UpdateUser()
	DeleteUser()
	GetUserByID(c *fiber.Ctx) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser() {

}
func (s *userService) GetUserByID(c *fiber.Ctx) (*models.User, error) {
	user, err := s.userRepo.GetAll(c.Context(), 10, 0)
	return user[0], err
}
func (s *userService) UpdateUser() {

}
func (s *userService) DeleteUser() {

}
