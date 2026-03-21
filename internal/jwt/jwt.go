package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/skibasu/auto-flow-api/internal/config"
)

type Claims struct {
	Sub   string   `json:"sub"`
	Roles []string `json:"role"`
	Type  string   `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(tokenType, userId string, roles []string, duration time.Duration) (string, error) {
	cfg := config.Load()
	secret := []byte(cfg.JWTSecret)

	claims := Claims{
		Sub:   userId,
		Roles: roles,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	cfg := config.Load()
	secret := []byte(cfg.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
