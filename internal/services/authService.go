package services

import (
	"errors"

	"time"

	"github.com/skibasu/auto-flow-api/internal/jwt"
	"github.com/skibasu/auto-flow-api/internal/repository"
)

const accessValidTime, refreshValidTime = 15 * time.Minute, 7 * 24 * time.Hour

type AuthService struct{ repo *repository.UserRepository }

func New(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Login(login, password string) (string, string, error) {
	user, err := s.repo.GetAuthDataByEmail(login)
	if err != nil {
		return "", "", err
	}

	// 🔐 sprawdzenie hasła
	// err = bcrypt.CompareHashAndPassword(
	// 	[]byte(user.Password),
	// 	[]byte(password),
	// )

	if password != user.Password {
		return "", "", errors.New("invalid credentials")
	} else if user.Password == password {
		access, _ := jwt.GenerateToken(user.Id, user.Roles, accessValidTime)
		refresh, _ := jwt.GenerateToken(user.Id, nil, refreshValidTime)

		return access, refresh, nil
	} else {
		return "", "", errors.New("invalid credentials")
	}

}

func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	userID := claims["sub"].(string)
	roles := claims["role"].([]string)

	access, _ := jwt.GenerateToken(userID, roles, accessValidTime)
	newRefresh, _ := jwt.GenerateToken(userID, nil, refreshValidTime)

	return access, newRefresh, nil
}
