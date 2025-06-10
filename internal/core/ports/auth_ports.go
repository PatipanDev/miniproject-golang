package ports

type AuthService interface {
	Login(email, password string) (string, error)
}
