package ports

import "test-backend/internal/core/domain"

//Secondery ports
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

//Primary ports
type UserService interface {
	RegisterUser(user *domain.User) error
}
