package port

import (
	"context"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/entity"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

type Repository interface {
	ListFunds(ctx context.Context) ([]entity.Fund, error)
	CreateFund(ctx context.Context, fund entity.Fund) (entity.Fund, error)
	UpdateFund(ctx context.Context, fund entity.Fund) (entity.Fund, error)
	GetFund(ctx context.Context, id vo.ID) (entity.Fund, error)

	ListInvestors(ctx context.Context) ([]entity.Investor, error)
	CreateInvestor(ctx context.Context, investor entity.Investor) (entity.Investor, error)
	GetInvestor(ctx context.Context, id vo.ID) (entity.Investor, error)

	ListInvestmentsByFund(ctx context.Context, fundID vo.ID) ([]entity.Investment, error)
	CreateInvestment(ctx context.Context, investment entity.Investment) (entity.Investment, error)
}
