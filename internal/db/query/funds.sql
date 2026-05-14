-- name: ListFunds :many
SELECT id, name, vintage_year, target_size_usd, status, created_at
FROM funds
ORDER BY created_at DESC;

-- name: CreateFund :one
INSERT INTO funds (name, vintage_year, target_size_usd, status)
VALUES ($1, $2, $3, $4)
RETURNING id, name, vintage_year, target_size_usd, status, created_at;

-- name: GetFund :one
SELECT id, name, vintage_year, target_size_usd, status, created_at
FROM funds
WHERE id = $1;

-- name: UpdateFund :one
UPDATE funds
SET name = $2,
    vintage_year = $3,
    target_size_usd = $4,
    status = $5
WHERE id = $1
RETURNING id, name, vintage_year, target_size_usd, status, created_at;

-- name: FundExists :one
SELECT EXISTS (
    SELECT 1
    FROM funds
    WHERE id = $1
);
