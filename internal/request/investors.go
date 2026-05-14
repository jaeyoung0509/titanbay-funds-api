package request

import (
	"strings"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"
)

type CreateInvestorRequest struct {
	Name         string `json:"name"`
	InvestorType string `json:"investor_type"`
	Email        string `json:"email"`
}

func (r *CreateInvestorRequest) Validate() map[string]string {
	fields := make(map[string]string)

	if strings.TrimSpace(r.Name) == "" {
		fields["name"] = "must be provided"
	}
	if _, err := enum.NewInvestorType(r.InvestorType); err != nil {
		fields["investor_type"] = err.Error()
	}
	if !IsEmail(r.Email) {
		fields["email"] = "must be a valid email address"
	}

	return fields
}

func (r CreateInvestorRequest) ToInput() (port.CreateInvestorInput, error) {
	investorType, err := enum.NewInvestorType(r.InvestorType)
	if err != nil {
		return port.CreateInvestorInput{}, domainerror.Validation("validation failed", map[string]string{
			"investor_type": err.Error(),
		})
	}

	email, err := vo.NewEmail(r.Email)
	if err != nil {
		return port.CreateInvestorInput{}, domainerror.Validation("validation failed", map[string]string{
			"email": "must be a valid email address",
		})
	}

	return port.CreateInvestorInput{
		Name:         strings.TrimSpace(r.Name),
		InvestorType: investorType,
		Email:        email,
	}, nil
}
