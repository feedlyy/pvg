package helper

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func GenerateActivationCodes() int {
	// Set seed for random number generator
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(8999) + 1000
}
