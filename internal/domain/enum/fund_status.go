package enum

import (
	"fmt"
	"strings"
)

type FundStatus string

const (
	FundStatusFundraising FundStatus = "Fundraising"
	FundStatusInvesting   FundStatus = "Investing"
	FundStatusClosed      FundStatus = "Closed"
)

func (s FundStatus) String() string {
	return string(s)
}

func (s FundStatus) Valid() bool {
	switch s {
	case FundStatusFundraising, FundStatusInvesting, FundStatusClosed:
		return true
	default:
		return false
	}
}

func FundStatusValues() []FundStatus {
	return []FundStatus{
		FundStatusFundraising,
		FundStatusInvesting,
		FundStatusClosed,
	}
}

func NewFundStatus(raw string) (FundStatus, error) {
	status := FundStatus(strings.TrimSpace(raw))
	if !status.Valid() {
		return "", fmt.Errorf("must be one of Fundraising, Investing, Closed")
	}
	return status, nil
}

