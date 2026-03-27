package password

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password *string) (string, error) {
	if password == nil || *password == "" {
		return "", nil
	}

	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(*password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	result := string(hashed)
	return result, nil
}
