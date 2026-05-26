-- name: ListInvestmentsByFund :many
SELECT id, investor_id, fund_id, amount_usd, investment_date
FROM investments
WHERE fund_id = $1
ORDER BY investment_date DESC, id DESC;

-- name: CreateInvestment :one
INSERT INTO investments (investor_id, fund_id, amount_usd, investment_date)
VALUES ($1, $2, $3, $4)
RETURNING id, investor_id, fund_id, amount_usd, investment_date;
