package password

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultBcryptCost = 12
)

var (
	ErrPasswordMismatch = errors.New("password does not match")
)

func HashPassword(password string) (string, string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", "", fmt.Errorf("error generating salt: %w", err)
	}

	saltStr := base64.StdEncoding.EncodeToString(salt)

	saltedPassword := password + saltStr

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), DefaultBcryptCost)
	if err != nil {
		return "", "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hashedBytes), saltStr, nil
}

func VerifyPassword(password, hash, salt string) error {
	saltedPassword := password + salt

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return fmt.Errorf("error verifying password: %w", err)
	}

	return nil
}