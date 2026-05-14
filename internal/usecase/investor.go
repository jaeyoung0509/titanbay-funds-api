package usecase

import (
	"context"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/entity"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"
)

type InvestorService interface {
	ListInvestors(ctx context.Context) ([]entity.Investor, error)
	CreateInvestor(ctx context.Context, input port.CreateInvestorInput) (entity.Investor, error)
	GetInvestor(ctx context.Context, id vo.ID) (entity.Investor, error)
}

type investorService struct {
	repo port.Repository
}

func NewInvestorService(repo port.Repository) InvestorService {
	return &investorService{repo: repo}
}

func (s *investorService) ListInvestors(ctx context.Context) ([]entity.Investor, error) {
	return s.repo.ListInvestors(ctx)
}

func (s *investorService) CreateInvestor(ctx context.Context, input port.CreateInvestorInput) (entity.Investor, error) {
	investor, err := entity.NewInvestor(input.Name, input.InvestorType, input.Email)
	if err != nil {
		return entity.Investor{}, err
	}

	return s.repo.CreateInvestor(ctx, investor)
}

func (s *investorService) GetInvestor(ctx context.Context, id vo.ID) (entity.Investor, error) {
	return s.repo.GetInvestor(ctx, id)
}
