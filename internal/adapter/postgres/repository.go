package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/db/sqlc"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/entity"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/enum"
	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/port"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	queries *sqlc.Queries
}

func New(db sqlc.DBTX) *Repository {
	return &Repository{queries: sqlc.New(db)}
}

var _ port.Repository = (*Repository)(nil)

func (r *Repository) ListFunds(ctx context.Context) ([]entity.Fund, error) {
	rows, err := r.queries.ListFunds(ctx)
	if err != nil {
		return nil, domainerror.Internal(err)
	}

	items := make([]entity.Fund, 0, len(rows))
	for _, row := range rows {
		item, err := mapFundRow(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) CreateFund(ctx context.Context, fund entity.Fund) (entity.Fund, error) {
	row, err := r.queries.CreateFund(ctx, sqlc.CreateFundParams{
		Name:          fund.Name,
		VintageYear:   int32(fund.VintageYear),
		TargetSizeUSD: fund.TargetSizeUSD.Decimal,
		Status:        fund.Status.String(),
	})
	if err != nil {
		return entity.Fund{}, mapDBError(err, "fund")
	}

	return mapFundRow(row)
}

func (r *Repository) UpdateFund(ctx context.Context, fund entity.Fund) (entity.Fund, error) {
	row, err := r.queries.UpdateFund(ctx, sqlc.UpdateFundParams{
		ID:            fund.ID.UUID,
		Name:          fund.Name,
		VintageYear:   int32(fund.VintageYear),
		TargetSizeUSD: fund.TargetSizeUSD.Decimal,
		Status:        fund.Status.String(),
	})
	if err != nil {
		return entity.Fund{}, mapDBError(err, "fund")
	}

	return mapFundRow(row)
}

func (r *Repository) GetFund(ctx context.Context, id vo.ID) (entity.Fund, error) {
	row, err := r.queries.GetFund(ctx, id.UUID)
	if err != nil {
		return entity.Fund{}, mapDBError(err, "fund")
	}
	return mapFundRow(row)
}

func (r *Repository) ListInvestors(ctx context.Context) ([]entity.Investor, error) {
	rows, err := r.queries.ListInvestors(ctx)
	if err != nil {
		return nil, domainerror.Internal(err)
	}

	items := make([]entity.Investor, 0, len(rows))
	for _, row := range rows {
		item, err := mapInvestorRow(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) CreateInvestor(ctx context.Context, investor entity.Investor) (entity.Investor, error) {
	row, err := r.queries.CreateInvestor(ctx, sqlc.CreateInvestorParams{
		Name:         investor.Name,
		InvestorType: investor.InvestorType.String(),
		Email:        investor.Email.String(),
	})
	if err != nil {
		return entity.Investor{}, mapDBError(err, "investor")
	}

	return mapInvestorRow(row)
}

func (r *Repository) GetInvestor(ctx context.Context, id vo.ID) (entity.Investor, error) {
	row, err := r.queries.GetInvestor(ctx, id.UUID)
	if err != nil {
		return entity.Investor{}, mapDBError(err, "investor")
	}
	return mapInvestorRow(row)
}

func (r *Repository) ListInvestmentsByFund(ctx context.Context, fundID vo.ID) ([]entity.Investment, error) {
	rows, err := r.queries.ListInvestmentsByFund(ctx, fundID.UUID)
	if err != nil {
		return nil, domainerror.Internal(err)
	}

	items := make([]entity.Investment, 0, len(rows))
	for _, row := range rows {
		item, err := mapInvestmentRow(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) CreateInvestment(ctx context.Context, investment entity.Investment) (entity.Investment, error) {
	row, err := r.queries.CreateInvestment(ctx, sqlc.CreateInvestmentParams{
		InvestorID:     investment.InvestorID.UUID,
		FundID:         investment.FundID.UUID,
		AmountUSD:      investment.AmountUSD.Decimal,
		InvestmentDate: investment.InvestmentDate.Time,
	})
	if err != nil {
		return entity.Investment{}, mapDBError(err, "investment")
	}

	return mapInvestmentRow(row)
}

func mapFundRow(row sqlc.Fund) (entity.Fund, error) {
	status, err := enum.NewFundStatus(row.Status)
	if err != nil {
		return entity.Fund{}, domainerror.Internal(fmt.Errorf("invalid fund status from db: %w", err))
	}

	item, err := entity.NewFundWithID(
		vo.NewIDFromUUID(row.ID),
		row.Name,
		int(row.VintageYear),
		vo.NewMoney(row.TargetSizeUSD),
		status,
		vo.NewTimestamp(row.CreatedAt),
	)
	if err != nil {
		return entity.Fund{}, domainerror.Internal(fmt.Errorf("invalid fund row: %w", err))
	}
	return item, nil
}

func mapInvestorRow(row sqlc.Investor) (entity.Investor, error) {
	investorType, err := enum.NewInvestorType(row.InvestorType)
	if err != nil {
		return entity.Investor{}, domainerror.Internal(fmt.Errorf("invalid investor type from db: %w", err))
	}

	email, err := vo.NewEmail(row.Email)
	if err != nil {
		return entity.Investor{}, domainerror.Internal(fmt.Errorf("invalid email from db: %w", err))
	}

	item, err := entity.NewInvestorWithID(
		vo.NewIDFromUUID(row.ID),
		row.Name,
		investorType,
		email,
		vo.NewTimestamp(row.CreatedAt),
	)
	if err != nil {
		return entity.Investor{}, domainerror.Internal(fmt.Errorf("invalid investor row: %w", err))
	}
	return item, nil
}

func mapInvestmentRow(row sqlc.Investment) (entity.Investment, error) {
	date := vo.NewDate(row.InvestmentDate)

	item, err := entity.NewInvestmentWithID(
		vo.NewIDFromUUID(row.ID),
		vo.NewIDFromUUID(row.FundID),
		vo.NewIDFromUUID(row.InvestorID),
		vo.NewMoney(row.AmountUSD),
		date,
	)
	if err != nil {
		return entity.Investment{}, domainerror.Internal(fmt.Errorf("invalid investment row: %w", err))
	}
	return item, nil
}

func mapDBError(err error, resource string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		if resource == "" {
			resource = "resource"
		}
		return domainerror.NotFound(resource)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			if strings.Contains(strings.ToLower(pgErr.ConstraintName), "email") {
				return domainerror.Conflict("email already exists")
			}
			return domainerror.Conflict("resource already exists")
		case "23503":
			return domainerror.NotFound("related resource")
		case "23514":
			return domainerror.Validation("validation failed", map[string]string{
				"body": "failed database constraint",
			})
		}
	}

	return domainerror.Internal(err)
}
