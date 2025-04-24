package password

import (
	"errors"
	"regexp"
)

const (
	MinPasswordLength = 8
)

var (
	ErrPasswordTooShort = errors.New("password is too short")
	ErrPasswordTooWeak  = errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
)

type PasswordConfig struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireDigits    bool
	RequireSpecial   bool
}

func DefaultPasswordConfig() PasswordConfig {
	return PasswordConfig{
		MinLength:        MinPasswordLength,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireDigits:    true,
		RequireSpecial:   true,
	}
}

func ValidatePassword(password string, config PasswordConfig) error {
	if len(password) < config.MinLength {
		return ErrPasswordTooShort
	}

	if config.RequireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return ErrPasswordTooWeak
	}

	if config.RequireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return ErrPasswordTooWeak
	}

	if config.RequireDigits && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return ErrPasswordTooWeak
	}

	if config.RequireSpecial && !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		return ErrPasswordTooWeak
	}

	return nil
}

func CheckPasswordStrength(password string) int {
	strength := 0

	if len(password) >= 8 {
		strength++
	}
	if len(password) >= 12 {
		strength++
	}

	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		strength++
	}
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		strength++
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		strength++
	}
	if regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		strength++
	}

	if len(password) < 8 {
		strength = 0
	}

	return strength
}
