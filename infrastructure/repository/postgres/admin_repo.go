package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"myapp/domain"
)

// AdminRepo is the PostgreSQL implementation of domain.AdminRepository.
type AdminRepo struct {
	pool *pgxpool.Pool
}

// NewAdminRepo creates a new AdminRepo.
func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{pool: pool}
}

// GetByEmail returns the admin with the given email, or ErrAdminNotFound.
func (r *AdminRepo) GetByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, email, password_hash, created_at, updated_at FROM admins WHERE email = $1`, email)
	return scanAdmin(row.Scan)
}

// GetByID returns the admin with the given ID, or ErrAdminNotFound.
func (r *AdminRepo) GetByID(ctx context.Context, id domain.AdminID) (*domain.Admin, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, email, password_hash, created_at, updated_at FROM admins WHERE id = $1`, string(id))
	return scanAdmin(row.Scan)
}

func scanAdmin(scan func(dest ...any) error) (*domain.Admin, error) {
	var (
		idStr string
		a     domain.Admin
	)
	err := scan(&idStr, &a.Name, &a.Email, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAdminNotFound
	}
	if err != nil {
		return nil, err
	}
	a.ID = domain.AdminID(idStr)
	return &a, nil
}
