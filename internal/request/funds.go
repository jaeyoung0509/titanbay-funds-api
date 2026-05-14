package request

import (
	"strings"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"
)

type CreateFundRequest struct {
	Name          string       `json:"name"`
	VintageYear   int          `json:"vintage_year"`
	TargetSizeUSD vo.Money     `json:"target_size_usd"`
	Status        string       `json:"status"`
}

func (r *CreateFundRequest) Validate() map[string]string {
	fields := make(map[string]string)

	if strings.TrimSpace(r.Name) == "" {
		fields["name"] = "must be provided"
	}
	if r.VintageYear < 1900 || r.VintageYear > 2100 {
		fields["vintage_year"] = "must be between 1900 and 2100"
	}
	if !r.TargetSizeUSD.IsPositive() {
		fields["target_size_usd"] = "must be greater than 0"
	} else if !r.TargetSizeUSD.HasMaxTwoDecimalPlaces() {
		fields["target_size_usd"] = "must have at most 2 decimal places"
	}
	if _, err := enum.NewFundStatus(r.Status); err != nil {
		fields["status"] = err.Error()
	}

	return fields
}

func (r CreateFundRequest) ToInput() (port.CreateFundInput, error) {
	status, err := enum.NewFundStatus(r.Status)
	if err != nil {
		return port.CreateFundInput{}, domainerror.Validation("validation failed", map[string]string{
			"status": err.Error(),
		})
	}

	return port.CreateFundInput{
		Name:          strings.TrimSpace(r.Name),
		VintageYear:   r.VintageYear,
		TargetSizeUSD: r.TargetSizeUSD,
		Status:        status,
	}, nil
}

type UpdateFundRequest struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	VintageYear   int          `json:"vintage_year"`
	TargetSizeUSD vo.Money     `json:"target_size_usd"`
	Status        string       `json:"status"`
}

func (r *UpdateFundRequest) Validate() map[string]string {
	fields := make(map[string]string)

	if _, err := vo.ParseID(r.ID); err != nil {
		fields["id"] = "must be a valid UUID"
	}
	if strings.TrimSpace(r.Name) == "" {
		fields["name"] = "must be provided"
	}
	if r.VintageYear < 1900 || r.VintageYear > 2100 {
		fields["vintage_year"] = "must be between 1900 and 2100"
	}
	if !r.TargetSizeUSD.IsPositive() {
		fields["target_size_usd"] = "must be greater than 0"
	} else if !r.TargetSizeUSD.HasMaxTwoDecimalPlaces() {
		fields["target_size_usd"] = "must have at most 2 decimal places"
	}
	if _, err := enum.NewFundStatus(r.Status); err != nil {
		fields["status"] = err.Error()
	}

	return fields
}

func (r UpdateFundRequest) ToInput() (port.UpdateFundInput, error) {
	id, err := vo.ParseID(r.ID)
	if err != nil {
		return port.UpdateFundInput{}, domainerror.Validation("validation failed", map[string]string{
			"id": "must be a valid UUID",
		})
	}

	status, err := enum.NewFundStatus(r.Status)
	if err != nil {
		return port.UpdateFundInput{}, domainerror.Validation("validation failed", map[string]string{
			"status": err.Error(),
		})
	}

	return port.UpdateFundInput{
		ID:            id,
		Name:          strings.TrimSpace(r.Name),
		VintageYear:   r.VintageYear,
		TargetSizeUSD: r.TargetSizeUSD,
		Status:        status,
	}, nil
}
