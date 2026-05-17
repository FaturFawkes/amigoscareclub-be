CREATE TABLE events (
    slug                    TEXT PRIMARY KEY,
    title                   TEXT NOT NULL,
    date                    DATE NOT NULL,
    time                    TEXT NOT NULL,
    timezone                TEXT NOT NULL DEFAULT 'Asia/Jakarta',
    location                TEXT NOT NULL,
    distance_km             INT  NOT NULL,
    pace                    TEXT NOT NULL,
    registration_open       BOOLEAN NOT NULL DEFAULT true,
    coffee_options          JSONB NOT NULL DEFAULT '[]',
    payment_bank            TEXT NOT NULL DEFAULT '',
    payment_account_number  TEXT NOT NULL DEFAULT '',
    payment_account_name    TEXT NOT NULL DEFAULT '',
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
