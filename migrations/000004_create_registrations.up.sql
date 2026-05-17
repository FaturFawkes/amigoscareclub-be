CREATE TABLE registrations (
    id                  TEXT PRIMARY KEY,
    ticket_number       TEXT UNIQUE NOT NULL,
    event_slug          TEXT NOT NULL REFERENCES events(slug),
    name                TEXT NOT NULL,
    email               TEXT NOT NULL,
    phone               TEXT NOT NULL,
    age                 INT  NOT NULL CHECK (age BETWEEN 10 AND 100),
    coffee_choice       TEXT NOT NULL,
    status              registration_status NOT NULL DEFAULT 'pending_verification',
    payment_proof_url   TEXT,
    note                TEXT,
    registered_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    verified_at         TIMESTAMPTZ,
    ticket_sent_at      TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (event_slug, email)
);

CREATE INDEX idx_registrations_event_slug ON registrations(event_slug);
CREATE INDEX idx_registrations_status     ON registrations(status);
