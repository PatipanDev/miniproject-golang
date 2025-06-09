package ports

<<<<<<< HEAD
import "github.com/PatipanDev/mini-project-golang/internal/core/domain"

type UserRepository interface {
	Update(user *domain.User, id string) error
	Delete(id string) error
}

type UserService interface {
	UpdateUser(user *domain.User, id string) error
	deleteUser(id string) error
=======
import "test-backend/internal/core/domain"

//Secondery ports
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

//Primary ports
type UserService interface {
	RegisterUser(user *domain.User) error
>>>>>>> origin/dev-nueng
}
