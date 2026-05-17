package dto

// GetEventInput holds the slug to look up an event.
type GetEventInput struct {
	Slug string
}

// EventPaymentData matches the swagger payment sub-schema.
type EventPaymentData struct {
	Bank          string `json:"bank"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

// EventData matches the swagger EventResponse data schema.
type EventData struct {
	Slug             string           `json:"slug"`
	Title            string           `json:"title"`
	Date             string           `json:"date"`
	Time             string           `json:"time"`
	Timezone         string           `json:"timezone"`
	Location         string           `json:"location"`
	DistanceKm       int              `json:"distance_km"`
	Pace             string           `json:"pace"`
	RegistrationOpen bool             `json:"registration_open"`
	CoffeeOptions    []string         `json:"coffee_options"`
	Payment          EventPaymentData `json:"payment"`
}

// GetEventOutput wraps a single event response.
type GetEventOutput struct {
	Data EventData `json:"data"`
}
