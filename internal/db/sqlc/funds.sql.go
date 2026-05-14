package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const listFundsQuery = `
SELECT id, name, vintage_year, target_size_usd, status, created_at
FROM funds
ORDER BY created_at DESC
`

const createFundQuery = `
INSERT INTO funds (name, vintage_year, target_size_usd, status)
VALUES ($1, $2, $3, $4)
RETURNING id, name, vintage_year, target_size_usd, status, created_at
`

const getFundQuery = `
SELECT id, name, vintage_year, target_size_usd, status, created_at
FROM funds
WHERE id = $1
`

const updateFundQuery = `
UPDATE funds
SET name = $2,
    vintage_year = $3,
    target_size_usd = $4,
    status = $5
WHERE id = $1
RETURNING id, name, vintage_year, target_size_usd, status, created_at
`

const fundExistsQuery = `
SELECT EXISTS (
    SELECT 1
    FROM funds
    WHERE id = $1
)
`

func (q *Queries) ListFunds(ctx context.Context) ([]Fund, error) {
	rows, err := q.db.QueryContext(ctx, listFundsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Fund
	for rows.Next() {
		var item Fund
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.VintageYear,
			&item.TargetSizeUSD,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) CreateFund(ctx context.Context, arg CreateFundParams) (Fund, error) {
	row := q.db.QueryRowContext(ctx, createFundQuery, arg.Name, arg.VintageYear, arg.TargetSizeUSD, arg.Status)
	var item Fund
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.VintageYear,
		&item.TargetSizeUSD,
		&item.Status,
		&item.CreatedAt,
	)
	return item, err
}

func (q *Queries) GetFund(ctx context.Context, id uuid.UUID) (Fund, error) {
	row := q.db.QueryRowContext(ctx, getFundQuery, id)
	var item Fund
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.VintageYear,
		&item.TargetSizeUSD,
		&item.Status,
		&item.CreatedAt,
	)
	return item, err
}

func (q *Queries) UpdateFund(ctx context.Context, arg UpdateFundParams) (Fund, error) {
	row := q.db.QueryRowContext(ctx, updateFundQuery, arg.ID, arg.Name, arg.VintageYear, arg.TargetSizeUSD, arg.Status)
	var item Fund
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.VintageYear,
		&item.TargetSizeUSD,
		&item.Status,
		&item.CreatedAt,
	)
	return item, err
}

func (q *Queries) FundExists(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRowContext(ctx, fundExistsQuery, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
