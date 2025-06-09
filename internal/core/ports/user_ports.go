package ports

import "github.com/PatipanDev/mini-project-golang/internal/core/domain"

//Secondery ports
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User, id string) error
	Delete(id string) error
	FindUserById(id string) (*domain.User, error)
}

//Primary ports
type UserService interface {
	RegisterUser(user *domain.User) error
	UpdateUser(user *domain.User, id string) error
	DeleteUser(id string) error
	GetUserById(id string) (*domain.User, error)
}
