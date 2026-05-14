package entity

import (
	"strings"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

type Investor struct {
	ID           vo.ID
	Name         string
	InvestorType enum.InvestorType
	Email        vo.Email
	CreatedAt    vo.Timestamp
}

func NewInvestor(name string, investorType enum.InvestorType, email vo.Email) (Investor, error) {
	fields := validateInvestor(name, investorType, email)
	if len(fields) > 0 {
		return Investor{}, domainerror.Validation("validation failed", fields)
	}

	return Investor{
		Name:         strings.TrimSpace(name),
		InvestorType: investorType,
		Email:        email,
	}, nil
}

func NewInvestorWithID(id vo.ID, name string, investorType enum.InvestorType, email vo.Email, createdAt vo.Timestamp) (Investor, error) {
	if id.IsZero() {
		return Investor{}, domainerror.Validation("validation failed", map[string]string{
			"id": "must be a valid UUID",
		})
	}
	fields := validateInvestor(name, investorType, email)
	if len(fields) > 0 {
		return Investor{}, domainerror.Validation("validation failed", fields)
	}

	return Investor{
		ID:           id,
		Name:         strings.TrimSpace(name),
		InvestorType: investorType,
		Email:        email,
		CreatedAt:    createdAt,
	}, nil
}

func validateInvestor(name string, investorType enum.InvestorType, email vo.Email) map[string]string {
	fields := make(map[string]string)
	if strings.TrimSpace(name) == "" {
		fields["name"] = "must be provided"
	}
	if !investorType.Valid() {
		fields["investor_type"] = "must be one of Individual, Institution, Family Office"
	}
	if email.String() == "" {
		fields["email"] = "must be a valid email address"
	}
	return fields
}

