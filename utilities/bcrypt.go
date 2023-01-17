package utilities

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken() (string, error) {
	randomToken := make([]byte, 32)
	_, err := rand.Read(randomToken)

	authToken := base64.URLEncoding.EncodeToString(randomToken)

	return authToken, err
}
