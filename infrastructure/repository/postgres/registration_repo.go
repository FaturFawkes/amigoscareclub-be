package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"myapp/domain"
)

const registrationSelectCols = `
	id, ticket_number, event_slug, name, email, phone, age, coffee_choice,
	status, payment_proof_url, note, registered_at, verified_at, ticket_sent_at,
	created_at, updated_at`

// RegistrationRepo is the PostgreSQL implementation of domain.RegistrationRepository.
type RegistrationRepo struct {
	pool *pgxpool.Pool
}

// NewRegistrationRepo creates a new RegistrationRepo.
func NewRegistrationRepo(pool *pgxpool.Pool) *RegistrationRepo {
	return &RegistrationRepo{pool: pool}
}

// Save inserts a new registration record.
func (r *RegistrationRepo) Save(ctx context.Context, reg *domain.Registration) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO registrations (
			id, ticket_number, event_slug, name, email, phone, age, coffee_choice,
			status, payment_proof_url, note, registered_at, verified_at, ticket_sent_at,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		string(reg.ID), reg.TicketNumber, reg.EventSlug,
		reg.Runner.Name, reg.Runner.Email, reg.Runner.Phone, reg.Runner.Age, reg.Runner.CoffeeChoice,
		string(reg.Status),
		nullableString(reg.PaymentProofURL), nullableString(reg.Note),
		reg.RegisteredAt, reg.VerifiedAt, reg.TicketSentAt,
		reg.CreatedAt, reg.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicateRegistration
		}
		return err
	}
	return nil
}

// GetByID retrieves a registration by event slug and registration ID.
func (r *RegistrationRepo) GetByID(ctx context.Context, eventSlug string, id domain.RegistrationID) (*domain.Registration, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT`+registrationSelectCols+` FROM registrations WHERE event_slug=$1 AND id=$2`,
		eventSlug, string(id))
	return scanRegistration(row.Scan)
}

// FindByEventAndEmail returns the registration matching the event + email pair, or ErrRegistrationNotFound.
func (r *RegistrationRepo) FindByEventAndEmail(ctx context.Context, eventSlug, email string) (*domain.Registration, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT`+registrationSelectCols+` FROM registrations WHERE event_slug=$1 AND email=$2`,
		eventSlug, email)
	return scanRegistration(row.Scan)
}

// List returns a filtered, paginated slice of registrations and the total matching count.
func (r *RegistrationRepo) List(ctx context.Context, eventSlug string, filter domain.RegistrationFilter, page, perPage int) ([]*domain.Registration, int, error) {
	args := []any{eventSlug}
	where := "event_slug = $1"
	if filter.Status != nil {
		args = append(args, string(*filter.Status))
		where += fmt.Sprintf(" AND status = $%d", len(args))
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM registrations WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	args = append(args, perPage, offset)
	dataSQL := fmt.Sprintf(
		`SELECT`+registrationSelectCols+` FROM registrations WHERE %s ORDER BY registered_at DESC LIMIT $%d OFFSET $%d`,
		where, len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var regs []*domain.Registration
	for rows.Next() {
		reg, err := scanRegistration(rows.Scan)
		if err != nil {
			return nil, 0, err
		}
		regs = append(regs, reg)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	return regs, total, nil
}

// Delete removes a registration by ID.
func (r *RegistrationRepo) Delete(ctx context.Context, id domain.RegistrationID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM registrations WHERE id=$1`, string(id))
	return err
}

// Update persists changes to an existing registration (status, note, timestamps).
func (r *RegistrationRepo) Update(ctx context.Context, reg *domain.Registration) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE registrations SET
			status=$1, note=$2, verified_at=$3, ticket_sent_at=$4, updated_at=$5
		WHERE id=$6`,
		string(reg.Status), nullableString(reg.Note),
		reg.VerifiedAt, reg.TicketSentAt, reg.UpdatedAt,
		string(reg.ID),
	)
	return err
}

// scanRegistration is a generic helper that works for both pgx.Row and pgx.Rows.
func scanRegistration(scan func(dest ...any) error) (*domain.Registration, error) {
	var (
		idStr, ticketNumber, eventSlug                    string
		name, email, phone, coffeeChoice, statusStr       string
		age                                               int
		paymentProofURL, note                             *string
		registeredAt, createdAt, updatedAt                time.Time
		verifiedAt, ticketSentAt                          *time.Time
	)
	err := scan(
		&idStr, &ticketNumber, &eventSlug,
		&name, &email, &phone, &age, &coffeeChoice,
		&statusStr,
		&paymentProofURL, &note,
		&registeredAt, &verifiedAt, &ticketSentAt,
		&createdAt, &updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrRegistrationNotFound
	}
	if err != nil {
		return nil, err
	}

	reg := &domain.Registration{
		ID:           domain.RegistrationID(idStr),
		TicketNumber: ticketNumber,
		EventSlug:    eventSlug,
		Runner: domain.Runner{
			Name:         name,
			Email:        email,
			Phone:        phone,
			Age:          age,
			CoffeeChoice: coffeeChoice,
		},
		Status:       domain.RegistrationStatus(statusStr),
		RegisteredAt: registeredAt,
		VerifiedAt:   verifiedAt,
		TicketSentAt: ticketSentAt,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
	if paymentProofURL != nil {
		reg.PaymentProofURL = *paymentProofURL
	}
	if note != nil {
		reg.Note = *note
	}
	return reg, nil
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
