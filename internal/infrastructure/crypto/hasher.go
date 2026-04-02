package crypto

import (
	"task-api/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func (BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (BcryptHasher) Compare(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return domain.ErrInvalidCredentials
	}
	return nil
}
