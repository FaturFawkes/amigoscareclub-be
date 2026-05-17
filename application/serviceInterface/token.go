package serviceInterface

import (
	"context"
	"time"
)

// TokenClaims holds the parsed data from a JWT token.
type TokenClaims struct {
	Sub       string
	JTI       string
	ExpiresAt time.Time
}

// TokenService issues and parses JWT tokens for admin authentication.
type TokenService interface {
	Issue(ctx context.Context, adminID string) (token string, claims TokenClaims, err error)
	Parse(ctx context.Context, token string) (TokenClaims, error)
}
