-- name: ListInvestors :many
SELECT id, name, investor_type, email, created_at
FROM investors
ORDER BY created_at DESC;

-- name: CreateInvestor :one
INSERT INTO investors (name, investor_type, email)
VALUES ($1, $2, $3)
RETURNING id, name, investor_type, email, created_at;

-- name: GetInvestor :one
SELECT id, name, investor_type, email, created_at
FROM investors
WHERE id = $1;

-- name: InvestorExists :one
SELECT EXISTS (
    SELECT 1
    FROM investors
    WHERE id = $1
);
