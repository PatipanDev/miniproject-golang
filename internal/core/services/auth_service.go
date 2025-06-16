package services

import (
	"errors"
	"time"

	"github.com/PatipanDev/mini-project-golang/internal/core/ports"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      ports.UserRepository
	jwtSecret string
}

func NewAurhService(repo ports.UserRepository, jwtsecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtsecret}
}

func (s *AuthService) Login(email, password string) (string, error) {
	// Get user from email
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	claims := &jwt.MapClaims{
		"sup":      user.ID,
		"email":    user.Email,
		"role":     user.Roles,
		"fullname": user.FirstName + " " + user.LastName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
