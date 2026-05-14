package entity

import (
	"strings"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

type Fund struct {
	ID            vo.ID
	Name          string
	VintageYear   int
	TargetSizeUSD vo.Money
	Status        enum.FundStatus
	CreatedAt     vo.Timestamp
}

func NewFund(name string, vintageYear int, targetSizeUSD vo.Money, status enum.FundStatus) (Fund, error) {
	fields := validateFund(name, vintageYear, targetSizeUSD, status)
	if len(fields) > 0 {
		return Fund{}, domainerror.Validation("validation failed", fields)
	}

	return Fund{
		Name:          strings.TrimSpace(name),
		VintageYear:   vintageYear,
		TargetSizeUSD: targetSizeUSD,
		Status:        status,
	}, nil
}

func NewFundWithID(id vo.ID, name string, vintageYear int, targetSizeUSD vo.Money, status enum.FundStatus, createdAt vo.Timestamp) (Fund, error) {
	if id.IsZero() {
		return Fund{}, domainerror.Validation("validation failed", map[string]string{
			"id": "must be a valid UUID",
		})
	}
	fields := validateFund(name, vintageYear, targetSizeUSD, status)
	if len(fields) > 0 {
		return Fund{}, domainerror.Validation("validation failed", fields)
	}

	return Fund{
		ID:            id,
		Name:          strings.TrimSpace(name),
		VintageYear:   vintageYear,
		TargetSizeUSD: targetSizeUSD,
		Status:        status,
		CreatedAt:     createdAt,
	}, nil
}

func validateFund(name string, vintageYear int, targetSizeUSD vo.Money, status enum.FundStatus) map[string]string {
	fields := make(map[string]string)
	if strings.TrimSpace(name) == "" {
		fields["name"] = "must be provided"
	}
	if vintageYear < 1900 || vintageYear > 2100 {
		fields["vintage_year"] = "must be between 1900 and 2100"
	}
	if !targetSizeUSD.IsPositive() {
		fields["target_size_usd"] = "must be greater than 0"
	} else if !targetSizeUSD.HasMaxTwoDecimalPlaces() {
		fields["target_size_usd"] = "must have at most 2 decimal places"
	}
	if !status.Valid() {
		fields["status"] = "must be one of Fundraising, Investing, Closed"
	}
	return fields
}

