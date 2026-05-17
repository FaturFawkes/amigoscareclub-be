package domain

import "context"

// EventRepository defines persistence operations for events.
type EventRepository interface {
	GetBySlug(ctx context.Context, slug string) (*Event, error)
}
