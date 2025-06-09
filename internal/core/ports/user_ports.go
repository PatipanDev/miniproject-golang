package ports

import "github.com/PatipanDev/mini-project-golang/internal/core/domain"

type UserRepository interface {
	Update(user *domain.User, id string) error
	Delete(id string) error
}

type UserService interface {
	UpdateUser(user *domain.User, id string) error
	deleteUser(id string) error
}
