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
		access, err := jwt.GenerateToken("access", user.Id, user.Roles, accessValidTime)
		if err != nil {
			return "", "", err
		}
		refresh, err := jwt.GenerateToken("refresh", user.Id, nil, refreshValidTime)
		if err != nil {
			return "", "", err
		}

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
	if claims.Type != "refresh" {
		return "", "", errors.New("invalid token type")
	}
	userID := claims.Sub
	roles := claims.Roles

	access, err := jwt.GenerateToken("access", userID, roles, accessValidTime)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := jwt.GenerateToken("refresh", userID, roles, refreshValidTime)
	if err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}
