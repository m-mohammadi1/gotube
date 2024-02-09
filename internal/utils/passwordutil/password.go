package passwordutil

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func PasswordMatches(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
