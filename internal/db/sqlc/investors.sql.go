package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const listInvestorsQuery = `
SELECT id, name, investor_type, email, created_at
FROM investors
ORDER BY created_at DESC
`

const createInvestorQuery = `
INSERT INTO investors (name, investor_type, email)
VALUES ($1, $2, $3)
RETURNING id, name, investor_type, email, created_at
`

const getInvestorQuery = `
SELECT id, name, investor_type, email, created_at
FROM investors
WHERE id = $1
`

const investorExistsQuery = `
SELECT EXISTS (
    SELECT 1
    FROM investors
    WHERE id = $1
)
`

func (q *Queries) ListInvestors(ctx context.Context) ([]Investor, error) {
	rows, err := q.db.QueryContext(ctx, listInvestorsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Investor
	for rows.Next() {
		var item Investor
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.InvestorType,
			&item.Email,
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

func (q *Queries) CreateInvestor(ctx context.Context, arg CreateInvestorParams) (Investor, error) {
	row := q.db.QueryRowContext(ctx, createInvestorQuery, arg.Name, arg.InvestorType, arg.Email)
	var item Investor
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.InvestorType,
		&item.Email,
		&item.CreatedAt,
	)
	return item, err
}

func (q *Queries) GetInvestor(ctx context.Context, id uuid.UUID) (Investor, error) {
	row := q.db.QueryRowContext(ctx, getInvestorQuery, id)
	var item Investor
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.InvestorType,
		&item.Email,
		&item.CreatedAt,
	)
	return item, err
}

func (q *Queries) InvestorExists(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRowContext(ctx, investorExistsQuery, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
