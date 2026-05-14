package enum

import (
	"fmt"
	"strings"
)

type InvestorType string

const (
	InvestorTypeIndividual   InvestorType = "Individual"
	InvestorTypeInstitution  InvestorType = "Institution"
	InvestorTypeFamilyOffice InvestorType = "Family Office"
)

func (t InvestorType) String() string {
	return string(t)
}

func (t InvestorType) Valid() bool {
	switch t {
	case InvestorTypeIndividual, InvestorTypeInstitution, InvestorTypeFamilyOffice:
		return true
	default:
		return false
	}
}

func InvestorTypeValues() []InvestorType {
	return []InvestorType{
		InvestorTypeIndividual,
		InvestorTypeInstitution,
		InvestorTypeFamilyOffice,
	}
}

func NewInvestorType(raw string) (InvestorType, error) {
	value := InvestorType(strings.TrimSpace(raw))
	if !value.Valid() {
		return "", fmt.Errorf("must be one of Individual, Institution, Family Office")
	}
	return value, nil
}

