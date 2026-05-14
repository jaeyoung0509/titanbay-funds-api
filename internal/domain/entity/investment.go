package entity

import (
	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

type Investment struct {
	ID             vo.ID
	InvestorID     vo.ID
	FundID         vo.ID
	AmountUSD      vo.Money
	InvestmentDate vo.Date
}

func NewInvestment(fundID vo.ID, investorID vo.ID, amountUSD vo.Money, investmentDate vo.Date) (Investment, error) {
	fields := validateInvestment(fundID, investorID, amountUSD, investmentDate)
	if len(fields) > 0 {
		return Investment{}, domainerror.Validation("validation failed", fields)
	}

	return Investment{
		FundID:         fundID,
		InvestorID:     investorID,
		AmountUSD:      amountUSD,
		InvestmentDate: investmentDate,
	}, nil
}

func NewInvestmentWithID(id vo.ID, fundID vo.ID, investorID vo.ID, amountUSD vo.Money, investmentDate vo.Date) (Investment, error) {
	if id.IsZero() {
		return Investment{}, domainerror.Validation("validation failed", map[string]string{
			"id": "must be a valid UUID",
		})
	}
	item, err := NewInvestment(fundID, investorID, amountUSD, investmentDate)
	if err != nil {
		return Investment{}, err
	}
	item.ID = id
	return item, nil
}

func validateInvestment(fundID vo.ID, investorID vo.ID, amountUSD vo.Money, investmentDate vo.Date) map[string]string {
	fields := make(map[string]string)
	if fundID.IsZero() {
		fields["fund_id"] = "must be a valid UUID"
	}
	if investorID.IsZero() {
		fields["investor_id"] = "must be a valid UUID"
	}
	if !amountUSD.IsPositive() {
		fields["amount_usd"] = "must be greater than 0"
	} else if !amountUSD.HasMaxTwoDecimalPlaces() {
		fields["amount_usd"] = "must have at most 2 decimal places"
	}
	if investmentDate.IsZero() {
		fields["investment_date"] = "must be a valid date"
	}
	return fields
}

