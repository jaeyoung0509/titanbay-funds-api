INSERT INTO funds (id, name, vintage_year, target_size_usd, status, created_at)
VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    'Titanbay Growth Fund I',
    2024,
    250000000.00,
    'Fundraising',
    '2024-01-15T10:30:00Z'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO investors (id, name, investor_type, email, created_at)
VALUES (
    '770e8400-e29b-41d4-a716-446655440002',
    'Goldman Sachs Asset Management',
    'Institution',
    'investments@gsam.com',
    '2024-02-10T09:15:00Z'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO investments (id, investor_id, fund_id, amount_usd, investment_date)
VALUES (
    '990e8400-e29b-41d4-a716-446655440004',
    '770e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440000',
    50000000.00,
    '2024-03-15'
)
ON CONFLICT (id) DO NOTHING;
