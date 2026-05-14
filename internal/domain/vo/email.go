package vo

import (
	"encoding/json"
	"errors"
	"net/mail"
	"strings"
)

var errInvalidAddress = errors.New("invalid email address")

type Email struct {
	Value string
}

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Email{}, errInvalidAddress
	}

	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return Email{}, err
	}
	if addr.Address != raw {
		return Email{}, errInvalidAddress
	}

	return Email{Value: raw}, nil
}

func (e Email) String() string {
	return e.Value
}

func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Value)
}

func (e *Email) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	value, err := NewEmail(raw)
	if err != nil {
		return err
	}
	*e = value
	return nil
}
