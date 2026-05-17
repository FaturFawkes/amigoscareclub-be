package idgen

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"myapp/infrastructure/config"
)

// Generator produces ULIDs for registrations/admins and sequential ticket numbers.
type Generator struct {
	pool   *pgxpool.Pool
	prefix string
}

// New creates a Generator using the config ticket prefix and a DB pool for the sequence.
func New(pool *pgxpool.Pool, cfg config.Config) *Generator {
	return &Generator{pool: pool, prefix: cfg.TicketPrefix}
}

// NewRegistrationID returns a new ULID prefixed with "reg_".
func (g *Generator) NewRegistrationID(_ context.Context) (string, error) {
	return "reg_" + ulid.Make().String(), nil
}

// NewAdminID returns a new ULID prefixed with "adm_".
func (g *Generator) NewAdminID(_ context.Context) (string, error) {
	return "adm_" + ulid.Make().String(), nil
}

// NewTicketNumber generates a sequential ticket number using a Postgres sequence.
func (g *Generator) NewTicketNumber(ctx context.Context) (string, error) {
	var seq int64
	if err := g.pool.QueryRow(ctx, "SELECT nextval('ticket_number_seq')").Scan(&seq); err != nil {
		return "", fmt.Errorf("nextval ticket_number_seq: %w", err)
	}
	return fmt.Sprintf("%s-%04d", g.prefix, seq), nil
}
