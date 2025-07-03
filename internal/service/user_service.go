package service

import (

	"go-backend-todo/internal/repository"

)

type UserService interface {
	CreateUser()
	UpdateUser()
	DeleteUser()
	GetUserByID()
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(){
	
}
func (s *userService) GetUserByID() {
	
}
func (s *userService) UpdateUser(){
	
}
func (s *userService) DeleteUser(){
	
}