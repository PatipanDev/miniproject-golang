package ports

import "github.com/PatipanDev/mini-project-golang/internal/core/domain"

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
}

type UserService interface {
	RegisterUser(user *domain.User) error
}
