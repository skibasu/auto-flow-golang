package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/skibasu/auto-flow-api/internal/consts"
)

func GenerateToken(userID string, roles []string, duration time.Duration) (string, error) {
	secret := []byte(consts.JWT_SECRET)
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
	}
	if roles != nil {
		claims["role"] = roles
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	secret := []byte(consts.JWT_SECRET)
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
