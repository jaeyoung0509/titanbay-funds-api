package enum

import "testing"

func TestFundStatus(t *testing.T) {
	status, err := NewFundStatus("Fundraising")
	if err != nil {
		t.Fatalf("new fund status: %v", err)
	}
	if status.String() != "Fundraising" {
		t.Fatalf("status = %q", status.String())
	}
	if !status.Valid() {
		t.Fatal("expected status to be valid")
	}

	if _, err := NewFundStatus("Draft"); err == nil {
		t.Fatal("expected invalid fund status to fail")
	}
}

func TestFundStatusValues(t *testing.T) {
	values := FundStatusValues()
	if len(values) != 3 {
		t.Fatalf("len = %d", len(values))
	}
}
