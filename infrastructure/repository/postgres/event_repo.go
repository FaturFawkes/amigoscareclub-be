package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"myapp/domain"
)

// EventRepo is the PostgreSQL implementation of domain.EventRepository.
type EventRepo struct {
	pool *pgxpool.Pool
}

// NewEventRepo creates a new EventRepo.
func NewEventRepo(pool *pgxpool.Pool) *EventRepo {
	return &EventRepo{pool: pool}
}

// GetBySlug returns the event matching the given slug, or ErrEventNotFound.
func (r *EventRepo) GetBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT slug, title, date, time, timezone, location, distance_km, pace,
		       registration_open, coffee_options,
		       payment_bank, payment_account_number, payment_account_name,
		       created_at, updated_at
		FROM events WHERE slug = $1`, slug)

	var (
		e          domain.Event
		dateVal    time.Time
		coffeeJSON []byte
	)
	err := row.Scan(
		&e.Slug, &e.Title, &dateVal, &e.Time, &e.Timezone, &e.Location,
		&e.DistanceKm, &e.Pace, &e.RegistrationOpen, &coffeeJSON,
		&e.Payment.Bank, &e.Payment.AccountNumber, &e.Payment.AccountName,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}

	e.Date = dateVal
	if err := json.Unmarshal(coffeeJSON, &e.CoffeeOptions); err != nil {
		return nil, err
	}
	return &e, nil
}
