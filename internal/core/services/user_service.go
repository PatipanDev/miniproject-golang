package services

import (
	"errors"
	"time"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/google/uuid"
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

	newUser := &domain.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     user.Email,
		Username:  user.Username,
		Password:  string(hassPassword),
		Status:    domain.USER_STATUS_ACTIVE,
		Roles:     make([]domain.Role, 0),
	}

	for _, r := range user.Roles {
		newUser.Roles = append(newUser.Roles, domain.Role{Name: domain.USER_ROLE(r.Name)})
	}

	return s.repo.Create(newUser)
}

func (s *UserServiceImp) UpdateUser(user *domain.User, id string) error {
	err := s.repo.Update(user, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImp) DeleteUser(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImp) GetUserById(id string) (*domain.User, error) {
	user, err := s.repo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
