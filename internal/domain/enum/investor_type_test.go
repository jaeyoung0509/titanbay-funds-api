package enum

import "testing"

func TestInvestorType(t *testing.T) {
	value, err := NewInvestorType("Institution")
	if err != nil {
		t.Fatalf("new investor type: %v", err)
	}
	if value.String() != "Institution" {
		t.Fatalf("value = %q", value.String())
	}
	if !value.Valid() {
		t.Fatal("expected investor type to be valid")
	}

	if _, err := NewInvestorType("Retail"); err == nil {
		t.Fatal("expected invalid investor type to fail")
	}
}

func TestInvestorTypeValues(t *testing.T) {
	values := InvestorTypeValues()
	if len(values) != 3 {
		t.Fatalf("len = %d", len(values))
	}
}
