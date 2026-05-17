package domain

import (
	"context"
	"time"
)

// TokenRepository manages JWT token revocation (blacklist).
type TokenRepository interface {
	Revoke(ctx context.Context, jti string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
}
