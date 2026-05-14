package handler

import (
	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
)

func validationError(field string) error {
	return domainerror.Validation("validation failed", map[string]string{
		field: "must be a valid UUID",
	})
}
