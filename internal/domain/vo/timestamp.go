package vo

import (
	"encoding/json"
	"time"
)

type Timestamp struct {
	Time time.Time
}

func NewTimestamp(value time.Time) Timestamp {
	return Timestamp{Time: value.UTC().Truncate(time.Second)}
}

func (t Timestamp) String() string {
	return t.Time.UTC().Format(time.RFC3339)
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return err
	}
	*t = NewTimestamp(parsed)
	return nil
}

