package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password []byte) ([]byte, error)
	Verify(password []byte, hash []byte) (bool, error)
}

var (
	defaultCost = 15
)

type BcryptHasher struct {
}

func (hasher *BcryptHasher) Hash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, defaultCost)
}

func (hasher *BcryptHasher) Verify(password []byte, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err == nil {
		return true, nil
	}

	return false, err
}
