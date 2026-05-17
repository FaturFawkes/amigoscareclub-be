package serviceInterface

import "context"

// PasswordHasher hashes and compares bcrypt passwords.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, hash, password string) error
}
