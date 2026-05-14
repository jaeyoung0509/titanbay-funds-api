package response

import (
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/entity"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

type Fund struct {
	ID            vo.ID           `json:"id"`
	Name          string          `json:"name"`
	VintageYear   int             `json:"vintage_year"`
	TargetSizeUSD vo.Money        `json:"target_size_usd"`
	Status        enum.FundStatus `json:"status"`
	CreatedAt     vo.Timestamp    `json:"created_at"`
}

func NewFund(item entity.Fund) Fund {
	return Fund{
		ID:            item.ID,
		Name:          item.Name,
		VintageYear:   item.VintageYear,
		TargetSizeUSD: item.TargetSizeUSD,
		Status:        item.Status,
		CreatedAt:     item.CreatedAt,
	}
}

type Investor struct {
	ID           vo.ID             `json:"id"`
	Name         string            `json:"name"`
	InvestorType enum.InvestorType `json:"investor_type"`
	Email        vo.Email          `json:"email"`
	CreatedAt    vo.Timestamp      `json:"created_at"`
}

func NewInvestor(item entity.Investor) Investor {
	return Investor{
		ID:           item.ID,
		Name:         item.Name,
		InvestorType: item.InvestorType,
		Email:        item.Email,
		CreatedAt:    item.CreatedAt,
	}
}

type Investment struct {
	ID             vo.ID        `json:"id"`
	InvestorID     vo.ID        `json:"investor_id"`
	FundID         vo.ID        `json:"fund_id"`
	AmountUSD      vo.Money     `json:"amount_usd"`
	InvestmentDate vo.Date      `json:"investment_date"`
}

func NewInvestment(item entity.Investment) Investment {
	return Investment{
		ID:             item.ID,
		InvestorID:     item.InvestorID,
		FundID:         item.FundID,
		AmountUSD:      item.AmountUSD,
		InvestmentDate: item.InvestmentDate,
	}
}

type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func NewErrorEnvelope(code, message string, fields map[string]string) ErrorEnvelope {
	return ErrorEnvelope{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	}
}
