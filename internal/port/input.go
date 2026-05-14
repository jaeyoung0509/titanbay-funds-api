package port

import (
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

type CreateFundInput struct {
	Name          string
	VintageYear   int
	TargetSizeUSD vo.Money
	Status        enum.FundStatus
}

type UpdateFundInput struct {
	ID            vo.ID
	Name          string
	VintageYear   int
	TargetSizeUSD vo.Money
	Status        enum.FundStatus
}

type CreateInvestorInput struct {
	Name         string
	InvestorType enum.InvestorType
	Email        vo.Email
}

type CreateInvestmentInput struct {
	FundID         vo.ID
	InvestorID     vo.ID
	AmountUSD      vo.Money
	InvestmentDate vo.Date
}
