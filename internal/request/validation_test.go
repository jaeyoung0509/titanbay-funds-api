package request

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
)

func TestCreateFundValidation(t *testing.T) {
	value, err := decimal.NewFromString("100.00")
	if err != nil {
		t.Fatalf("decimal parse: %v", err)
	}

	req := CreateFundRequest{
		Name:          "  ",
		VintageYear:   1800,
		TargetSizeUSD: vo.NewMoney(value),
		Status:        "Draft",
	}

	fields := req.Validate()
	if fields["name"] == "" {
		t.Fatal("expected name validation error")
	}
	if fields["vintage_year"] == "" {
		t.Fatal("expected vintage_year validation error")
	}
	if fields["status"] == "" {
		t.Fatal("expected status validation error")
	}
}

func TestCreateInvestorValidation(t *testing.T) {
	req := CreateInvestorRequest{
		Name:         "",
		InvestorType: "Retail",
		Email:        "not-an-email",
	}

	fields := req.Validate()
	if fields["name"] == "" {
		t.Fatal("expected name validation error")
	}
	if fields["investor_type"] == "" {
		t.Fatal("expected investor_type validation error")
	}
	if fields["email"] == "" {
		t.Fatal("expected email validation error")
	}
}

func TestCreateInvestmentValidation(t *testing.T) {
	invalidAmount, err := decimal.NewFromString("1.234")
	if err != nil {
		t.Fatalf("decimal parse: %v", err)
	}

	req := CreateInvestmentRequest{
		InvestorID:     "bad",
		AmountUSD:      vo.NewMoney(invalidAmount),
		InvestmentDate: "",
	}

	fields := req.Validate()
	if fields["investor_id"] == "" {
		t.Fatal("expected investor_id validation error")
	}
	if fields["amount_usd"] == "" {
		t.Fatal("expected amount_usd validation error")
	}
	if fields["investment_date"] == "" {
		t.Fatal("expected investment_date validation error")
	}
}

func TestUpdateFundValidation(t *testing.T) {
	value, err := decimal.NewFromString("100.00")
	if err != nil {
		t.Fatalf("decimal parse: %v", err)
	}

	req := UpdateFundRequest{
		ID:            "bad",
		Name:          "",
		VintageYear:   2500,
		TargetSizeUSD: vo.NewMoney(value),
		Status:        "Draft",
	}

	fields := req.Validate()
	if fields["id"] == "" {
		t.Fatal("expected id validation error")
	}
	if fields["name"] == "" {
		t.Fatal("expected name validation error")
	}
	if fields["vintage_year"] == "" {
		t.Fatal("expected vintage_year validation error")
	}
	if fields["status"] == "" {
		t.Fatal("expected status validation error")
	}
}
