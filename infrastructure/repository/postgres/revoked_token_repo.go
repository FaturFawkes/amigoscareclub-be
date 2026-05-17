package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"myapp/domain"
)

// RevokedTokenRepo is the PostgreSQL implementation of domain.TokenRepository.
type RevokedTokenRepo struct {
	pool *pgxpool.Pool
}

// NewRevokedTokenRepo creates a new RevokedTokenRepo.
func NewRevokedTokenRepo(pool *pgxpool.Pool) *RevokedTokenRepo {
	return &RevokedTokenRepo{pool: pool}
}

// Revoke adds the JTI to the blacklist table.
func (r *RevokedTokenRepo) Revoke(ctx context.Context, jti string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO revoked_tokens (jti, expires_at) VALUES ($1, $2) ON CONFLICT (jti) DO NOTHING`,
		jti, expiresAt)
	return err
}

// IsRevoked returns true if the JTI is present in the blacklist and not yet expired.
func (r *RevokedTokenRepo) IsRevoked(ctx context.Context, jti string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM revoked_tokens WHERE jti = $1 AND expires_at > NOW()`, jti).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Ensure RevokedTokenRepo satisfies domain.TokenRepository at compile time.
var _ domain.TokenRepository = (*RevokedTokenRepo)(nil)
