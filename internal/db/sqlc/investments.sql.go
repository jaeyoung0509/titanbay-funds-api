package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const listInvestmentsByFundQuery = `
SELECT id, investor_id, fund_id, amount_usd, investment_date::text
FROM investments
WHERE fund_id = $1
ORDER BY investment_date DESC, id DESC
`

const createInvestmentQuery = `
INSERT INTO investments (investor_id, fund_id, amount_usd, investment_date)
VALUES ($1, $2, $3, $4)
RETURNING id, investor_id, fund_id, amount_usd, investment_date::text
`

func (q *Queries) ListInvestmentsByFund(ctx context.Context, fundID uuid.UUID) ([]Investment, error) {
	rows, err := q.db.QueryContext(ctx, listInvestmentsByFundQuery, fundID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Investment
	for rows.Next() {
		var item Investment
		if err := rows.Scan(
			&item.ID,
			&item.InvestorID,
			&item.FundID,
			&item.AmountUSD,
			&item.InvestmentDate,
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

func (q *Queries) CreateInvestment(ctx context.Context, arg CreateInvestmentParams) (Investment, error) {
	row := q.db.QueryRowContext(ctx, createInvestmentQuery, arg.InvestorID, arg.FundID, arg.AmountUSD, arg.InvestmentDate)
	var item Investment
	err := row.Scan(
		&item.ID,
		&item.InvestorID,
		&item.FundID,
		&item.AmountUSD,
		&item.InvestmentDate,
	)
	return item, err
}
