package vo

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

type ID struct {
	UUID uuid.UUID
}

func NewID() ID {
	return ID{UUID: uuid.New()}
}

func NewIDFromUUID(value uuid.UUID) ID {
	return ID{UUID: value}
}

func ParseID(raw string) (ID, error) {
	value, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ID{}, err
	}
	return ID{UUID: value}, nil
}

func (id ID) IsZero() bool {
	return id.UUID == uuid.Nil
}

func (id ID) String() string {
	return id.UUID.String()
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed, err := ParseID(raw)
	if err != nil {
		return err
	}
	*id = parsed
	return nil
}

