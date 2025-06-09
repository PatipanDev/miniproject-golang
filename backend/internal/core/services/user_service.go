package services

import "github.com/PatipanDev/mini-project-golang/internal/core/ports"

type userService struct {
	service ports.UserRepository 
}

func NewUserService (service ports.UserRepository) ports.UserService{
	return userService{service: service}
}

func (s *userService) RegisterUser(user *)