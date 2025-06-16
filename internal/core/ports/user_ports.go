package ports

import "github.com/PatipanDev/mini-project-golang/internal/core/domain"

//Secondery ports
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User, id string) error
	Delete(id string) error
	FindUserById(id string) (*domain.User, error)
	//GetProfile(id string) (*domain.User, error)
	Get() ([]domain.User, error)
	FindAll(offset int, limit int) ([]domain.User, error)
	Count() (int64, error)
	//filter
	FindUsers(filter *domain.UserFilter) ([]domain.User, int64, error)
	//update profile
	UpdateUserProfilePicURL(id string, url string) error
}

type FileStorageRepository interface {
	SaveFile(folderPath string, filename string, fileContent []byte) (string, error)
}

//Primary ports
type UserService interface {
	RegisterUser(user *domain.User) error
	UpdateUser(user *domain.User, id string) error
	DeleteUser(id string) error
	GetUserById(id string) (*domain.User, error)
	FindUsers() ([]domain.User, error)
	GetPaginationUsers(page int, limit int) (*domain.Pagination, error)
	GetUsers(filter *domain.UserFilter) ([]domain.ResUSerNoID, int64, error)
	//upload file image
	UploadProfilePicture(id string, file []byte, filename string) (string, error)
}
