package vo

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

type Money struct {
	Decimal decimal.Decimal
}

func NewMoney(value decimal.Decimal) Money {
	return Money{Decimal: value}
}

func ParseMoney(raw string) (Money, error) {
	value, err := decimal.NewFromString(strings.TrimSpace(raw))
	if err != nil {
		return Money{}, err
	}
	return NewMoney(value), nil
}

func (m Money) IsPositive() bool {
	return m.Decimal.GreaterThan(decimal.Zero)
}

func (m Money) HasMaxTwoDecimalPlaces() bool {
	return m.Decimal.Equal(m.Decimal.Round(2))
}

func (m Money) String() string {
	return m.Decimal.StringFixed(2)
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return errors.New("money value is required")
	}

	if trimmed[0] == '"' {
		var raw string
		if err := json.Unmarshal(trimmed, &raw); err != nil {
			return err
		}
		value, err := decimal.NewFromString(strings.TrimSpace(raw))
		if err != nil {
			return err
		}
		m.Decimal = value
		return nil
	}

	value, err := decimal.NewFromString(string(trimmed))
	if err != nil {
		return err
	}
	m.Decimal = value
	return nil
}

