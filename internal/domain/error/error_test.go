package domainerror

import (
	"errors"
	"testing"
)

func TestConstructors(t *testing.T) {
	validation := Validation("validation failed", map[string]string{"status": "invalid"})
	if validation.Kind != KindValidation {
		t.Fatalf("kind = %s", validation.Kind)
	}

	notFound := NotFound("fund")
	if notFound.Kind != KindNotFound {
		t.Fatalf("kind = %s", notFound.Kind)
	}

	conflict := Conflict("duplicate")
	if conflict.Kind != KindConflict {
		t.Fatalf("kind = %s", conflict.Kind)
	}

	internal := Internal(errors.New("boom"))
	if internal.Kind != KindInternal {
		t.Fatalf("kind = %s", internal.Kind)
	}
}

func TestAs(t *testing.T) {
	original := Validation("validation failed", nil)
	if got := As(original); got != original {
		t.Fatal("expected As to return original domain error")
	}

	wrapped := As(errors.New("boom"))
	if wrapped.Kind != KindInternal {
		t.Fatalf("kind = %s", wrapped.Kind)
	}
}
