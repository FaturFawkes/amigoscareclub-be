package domain

import "time"

// EventPayment holds bank transfer details for event registration fee.
type EventPayment struct {
	Bank          string
	AccountNumber string
	AccountName   string
}

// Event represents a running event that participants can register for.
type Event struct {
	Slug             string
	Title            string
	Date             time.Time
	Time             string
	Timezone         string
	Location         string
	DistanceKm       int
	Pace             string
	RegistrationOpen bool
	CoffeeOptions    []string
	Payment          EventPayment
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
