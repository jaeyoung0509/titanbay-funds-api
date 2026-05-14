package request

import (
	"strings"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"
)

type CreateInvestmentRequest struct {
	InvestorID     string       `json:"investor_id"`
	AmountUSD      vo.Money     `json:"amount_usd"`
	InvestmentDate string       `json:"investment_date"`
}

func (r *CreateInvestmentRequest) Validate() map[string]string {
	fields := make(map[string]string)

	if _, err := vo.ParseID(r.InvestorID); err != nil {
		fields["investor_id"] = "must be a valid UUID"
	}
	if !r.AmountUSD.IsPositive() {
		fields["amount_usd"] = "must be greater than 0"
	} else if !r.AmountUSD.HasMaxTwoDecimalPlaces() {
		fields["amount_usd"] = "must have at most 2 decimal places"
	}
	if _, err := vo.ParseDate(strings.TrimSpace(r.InvestmentDate)); err != nil {
		fields["investment_date"] = "must be a valid date"
	}

	return fields
}

func (r CreateInvestmentRequest) ToInput(fundID vo.ID) (port.CreateInvestmentInput, error) {
	investorID, err := vo.ParseID(r.InvestorID)
	if err != nil {
		return port.CreateInvestmentInput{}, domainerror.Validation("validation failed", map[string]string{
			"investor_id": "must be a valid UUID",
		})
	}

	investmentDate, err := vo.ParseDate(strings.TrimSpace(r.InvestmentDate))
	if err != nil {
		return port.CreateInvestmentInput{}, domainerror.Validation("validation failed", map[string]string{
			"investment_date": "must be a valid date",
		})
	}

	return port.CreateInvestmentInput{
		FundID:         fundID,
		InvestorID:     investorID,
		AmountUSD:      r.AmountUSD,
		InvestmentDate: investmentDate,
	}, nil
}
