package usecase

import (
	"context"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/entity"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"
)

type FundService interface {
	ListFunds(ctx context.Context) ([]entity.Fund, error)
	CreateFund(ctx context.Context, input port.CreateFundInput) (entity.Fund, error)
	UpdateFund(ctx context.Context, input port.UpdateFundInput) (entity.Fund, error)
	GetFund(ctx context.Context, id vo.ID) (entity.Fund, error)
}

type fundService struct {
	repo port.Repository
}

func NewFundService(repo port.Repository) FundService {
	return &fundService{repo: repo}
}

func (s *fundService) ListFunds(ctx context.Context) ([]entity.Fund, error) {
	return s.repo.ListFunds(ctx)
}

func (s *fundService) CreateFund(ctx context.Context, input port.CreateFundInput) (entity.Fund, error) {
	fund, err := entity.NewFund(input.Name, input.VintageYear, input.TargetSizeUSD, input.Status)
	if err != nil {
		return entity.Fund{}, err
	}

	return s.repo.CreateFund(ctx, fund)
}

func (s *fundService) UpdateFund(ctx context.Context, input port.UpdateFundInput) (entity.Fund, error) {
	if input.ID.IsZero() {
		return entity.Fund{}, domainerror.Validation("validation failed", map[string]string{
			"id": "must be a valid UUID",
		})
	}

	fund, err := entity.NewFundWithID(input.ID, input.Name, input.VintageYear, input.TargetSizeUSD, input.Status, vo.Timestamp{})
	if err != nil {
		return entity.Fund{}, err
	}

	return s.repo.UpdateFund(ctx, fund)
}

func (s *fundService) GetFund(ctx context.Context, id vo.ID) (entity.Fund, error) {
	return s.repo.GetFund(ctx, id)
}
