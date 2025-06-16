package ports

import "github.com/PatipanDev/mini-project-golang/internal/core/domain"

type RoleRepository interface {
	FindByName(name string, role *domain.Role) error
	Create(role *domain.Role) error
}
