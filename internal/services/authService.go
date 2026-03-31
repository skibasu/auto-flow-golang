package services

import (
	"errors"

	"time"

	"github.com/jackc/pgx/v5"
	"github.com/skibasu/auto-flow-api/internal/jwt"
	"golang.org/x/crypto/bcrypt"
)

const accessValidTime, refreshValidTime = 15 * time.Minute, 7 * 24 * time.Hour

func (s *Service) Login(login, password string) (string, string, error) {
	user, err := s.repo.GetAuthDataByEmail(login)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", errors.New("invalid credentials")
		}

		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	access, err := jwt.GenerateToken("access", user.Id, s.config.Secret, user.Roles, accessValidTime)
	if err != nil {
		return "", "", err
	}
	refresh, err := jwt.GenerateToken("refresh", user.Id, s.config.Secret, user.Roles, refreshValidTime)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil

}

func (s *Service) Refresh(refreshToken string) (string, string, error) {
	claims, err := jwt.ParseToken(refreshToken, s.config.Secret)
	if err != nil {
		return "", "", err
	}
	if claims.Type != "refresh" {
		return "", "", errors.New("invalid token type")
	}
	userID := claims.Sub
	roles := claims.Roles

	access, err := jwt.GenerateToken("access", userID, s.config.Secret, roles, accessValidTime)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := jwt.GenerateToken("refresh", userID, s.config.Secret, roles, refreshValidTime)
	if err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}
