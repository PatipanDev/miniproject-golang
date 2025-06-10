package services

import (
	"errors"
	"fmt"
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

func (s *UserServiceImp) FindUsers() ([]domain.User, error) {
	r, err := s.repo.Get()
	if err != nil {
		return nil, fmt.Errorf("error get users : %w", err)
	}
	return r, nil
}

func (s *UserServiceImp) SearchUser(qurey string) ([]domain.User, error) {
	return s.repo.SearchData(qurey)
}

func (s *UserServiceImp) GetPaginationUsers(page int, limit int) (*domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	total, err := s.repo.Count()
	if err != nil {
		return nil, err
	}

	users, err := s.repo.FindAll(offset, limit)
	if err != nil {
		return nil, err
	}

	totalPage := int((total + int64(limit) - 1) / int64(limit))

	return &domain.Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
		Data:      users,
	}, nil
}

func (s *UserServiceImp) GetUsers(filter *domain.UserFilter) ([]domain.User, int64, error) {
	users, total, err := s.repo.FindUsers(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error get users: %w", err)
	}
	return users, total, nil
}
