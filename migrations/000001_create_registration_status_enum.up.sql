CREATE TYPE registration_status AS ENUM (
    'pending_verification',
    'verified',
    'rejected',
    'ticket_sent'
);
