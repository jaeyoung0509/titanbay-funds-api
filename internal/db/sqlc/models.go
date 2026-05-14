package sqlc

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Fund struct {
	ID            uuid.UUID
	Name          string
	VintageYear   int32
	TargetSizeUSD decimal.Decimal
	Status        string
	CreatedAt     time.Time
}

type Investor struct {
	ID           uuid.UUID
	Name         string
	InvestorType string
	Email        string
	CreatedAt    time.Time
}

type Investment struct {
	ID             uuid.UUID
	InvestorID     uuid.UUID
	FundID         uuid.UUID
	AmountUSD      decimal.Decimal
	InvestmentDate string
}

type CreateFundParams struct {
	Name          string
	VintageYear   int32
	TargetSizeUSD decimal.Decimal
	Status        string
}

type UpdateFundParams struct {
	ID            uuid.UUID
	Name          string
	VintageYear   int32
	TargetSizeUSD decimal.Decimal
	Status        string
}

type CreateInvestorParams struct {
	Name         string
	InvestorType string
	Email        string
}

type CreateInvestmentParams struct {
	InvestorID     uuid.UUID
	FundID         uuid.UUID
	AmountUSD      decimal.Decimal
	InvestmentDate time.Time
}
