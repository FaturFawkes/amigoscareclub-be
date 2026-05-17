INSERT INTO events (
    slug, title, date, time, timezone, location, distance_km, pace,
    registration_open, coffee_options,
    payment_bank, payment_account_number, payment_account_name
) VALUES (
    '40-of-heart-rate-run',
    '40% OF HEART RATE RUN – VOL.2',
    '2026-05-24',
    '06:00',
    'Asia/Jakarta',
    'Melkkops Coffee & Eatry',
    5,
    'Every pace welcome',
    true,
    '["Americano","Cappuccino","Latte","Es Kopi Susu","Espresso"]',
    'BCA',
    '4061207427',
    'Nur Fatchurohman'
) ON CONFLICT (slug) DO NOTHING;
