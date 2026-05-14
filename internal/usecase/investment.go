package usecase

import (
	"context"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/entity"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"
)

type InvestmentService interface {
	ListInvestmentsByFund(ctx context.Context, fundID vo.ID) ([]entity.Investment, error)
	CreateInvestment(ctx context.Context, input port.CreateInvestmentInput) (entity.Investment, error)
}

type investmentService struct {
	repo port.Repository
}

func NewInvestmentService(repo port.Repository) InvestmentService {
	return &investmentService{repo: repo}
}

func (s *investmentService) ListInvestmentsByFund(ctx context.Context, fundID vo.ID) ([]entity.Investment, error) {
	if _, err := s.repo.GetFund(ctx, fundID); err != nil {
		return nil, err
	}

	return s.repo.ListInvestmentsByFund(ctx, fundID)
}

func (s *investmentService) CreateInvestment(ctx context.Context, input port.CreateInvestmentInput) (entity.Investment, error) {
	investment, err := entity.NewInvestment(input.FundID, input.InvestorID, input.AmountUSD, input.InvestmentDate)
	if err != nil {
		return entity.Investment{}, err
	}

	if _, err := s.repo.GetFund(ctx, input.FundID); err != nil {
		return entity.Investment{}, err
	}
	if _, err := s.repo.GetInvestor(ctx, input.InvestorID); err != nil {
		return entity.Investment{}, err
	}

	return s.repo.CreateInvestment(ctx, investment)
}
