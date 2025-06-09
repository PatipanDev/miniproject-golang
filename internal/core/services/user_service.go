package services

import (
	"errors"
	"test-backend/internal/core/domain"
	"test-backend/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImp struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserServiceImp{repo: repo}
}

func (s *UserServiceImp) RegisterUser(user *domain.User) error {
	eixsting, _ := s.repo.FindByEmail(user.Email)
	if eixsting != nil {
		return errors.New("email already in use")
	}

	hassPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hassPassword)
	user.Status = domain.USER_STATUS_ACTIVE

	return s.repo.Create(user)
}
