package vo

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestMoneyJSONRoundTrip(t *testing.T) {
	value, err := decimal.NewFromString("123.40")
	if err != nil {
		t.Fatalf("decimal parse: %v", err)
	}

	money := NewMoney(value)

	data, err := json.Marshal(money)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := string(data); got != "123.40" {
		t.Fatalf("marshal = %q, want %q", got, "123.40")
	}

	var decoded Money
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := decoded.String(); got != "123.40" {
		t.Fatalf("decoded = %q, want %q", got, "123.40")
	}
}

func TestMoneyValidation(t *testing.T) {
	positive, _ := decimal.NewFromString("1.23")
	invalidScale, _ := decimal.NewFromString("1.234")
	zero, _ := decimal.NewFromString("0")
	negative, _ := decimal.NewFromString("-1.23")

	if !NewMoney(positive).IsPositive() {
		t.Fatal("expected positive money to be valid")
	}
	if !NewMoney(positive).HasMaxTwoDecimalPlaces() {
		t.Fatal("expected 2dp value to be valid")
	}
	if NewMoney(invalidScale).HasMaxTwoDecimalPlaces() {
		t.Fatal("expected scale > 2 to be invalid")
	}
	if NewMoney(zero).IsPositive() {
		t.Fatal("expected zero to be invalid")
	}
	if NewMoney(negative).IsPositive() {
		t.Fatal("expected negative to be invalid")
	}
}
