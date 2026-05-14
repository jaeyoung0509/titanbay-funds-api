package vo

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateJSONRoundTrip(t *testing.T) {
	date, err := ParseDate("2024-09-22")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	data, err := json.Marshal(date)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := string(data); got != `"2024-09-22"` {
		t.Fatalf("marshal = %q, want %q", got, `"2024-09-22"`)
	}

	var decoded Date
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := decoded.String(); got != "2024-09-22" {
		t.Fatalf("decoded = %q, want %q", got, "2024-09-22")
	}
}

func TestTimestampJSON(t *testing.T) {
	ts := NewTimestamp(time.Date(2024, time.January, 15, 10, 30, 0, 123000000, time.FixedZone("BST", 3600)))
	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := string(data); got != `"2024-01-15T09:30:00Z"` {
		t.Fatalf("marshal = %q, want %q", got, `"2024-01-15T09:30:00Z"`)
	}
}
