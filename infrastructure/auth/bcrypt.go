package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements PasswordHasher using bcrypt.
type BcryptHasher struct{}

// NewBcryptHasher creates a BcryptHasher.
func NewBcryptHasher() *BcryptHasher { return &BcryptHasher{} }

// Hash generates a bcrypt hash from the plain-text password.
func (h *BcryptHasher) Hash(_ context.Context, password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare returns nil when the hash matches the password.
func (h *BcryptHasher) Compare(_ context.Context, hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
