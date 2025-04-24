package token

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	Alphabetic = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	URLSafe = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("error generating random bytes: %w", err)
	}
	return bytes, nil
}

func GenerateRandomString(length int, charset string) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("error generating random string: %w", err)
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

func GenerateRandomAlphanumeric(length int) (string, error) {
	return GenerateRandomString(length, Alphanumeric)
}

func GenerateVerificationToken() (string, error) {
	return GenerateRandomString(32, URLSafe)
}

func GenerateResetToken() (string, error) {
	return GenerateRandomString(64, URLSafe)
}

func GenerateRefreshToken() (string, error) {
	bytes, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateSessionID() (string, error) {
	bytes, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}