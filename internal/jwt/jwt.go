package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub   string   `json:"sub"`
	Roles []string `json:"role"`
	Type  string   `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(tokenType, userId, secret string, roles []string, duration time.Duration) (string, error) {

	s := []byte(secret)

	claims := Claims{
		Sub:   userId,
		Roles: roles,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s)
}

func ParseToken(tokenStr, secret string) (*Claims, error) {

	s := []byte(secret)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return s, nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
