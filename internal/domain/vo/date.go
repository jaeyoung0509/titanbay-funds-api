package vo

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"
)

const dateLayout = "2006-01-02"

type Date struct {
	Time time.Time
}

func NewDate(value time.Time) Date {
	utc := value.UTC()
	return Date{
		Time: time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC),
	}
}

func ParseDate(raw string) (Date, error) {
	value, err := time.Parse(dateLayout, raw)
	if err != nil {
		return Date{}, err
	}
	return NewDate(value), nil
}

func (d Date) IsZero() bool {
	return d.Time.IsZero()
}

func (d Date) String() string {
	return d.Time.UTC().Format(dateLayout)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return errors.New("date value is required")
	}

	var raw string
	if err := json.Unmarshal(trimmed, &raw); err != nil {
		return err
	}

	parsed, err := ParseDate(raw)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

