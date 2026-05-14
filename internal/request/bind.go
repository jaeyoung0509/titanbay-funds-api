package request

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gofiber/fiber/v3"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
)

type Validator interface {
	Validate() map[string]string
}

func BindAndValidate[T any](c fiber.Ctx) (*T, error) {
	body := bytes.TrimSpace(c.Body())
	if len(body) == 0 {
		return nil, domainerror.Validation("validation failed", map[string]string{
			"body": "request body is required",
		})
	}

	var payload T
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		return nil, domainerror.Validation("validation failed", map[string]string{
			"body": "invalid JSON",
		})
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, domainerror.Validation("validation failed", map[string]string{
			"body": "invalid JSON",
		})
	}

	if validator, ok := any(&payload).(Validator); ok {
		if fields := validator.Validate(); len(fields) > 0 {
			return nil, domainerror.Validation("validation failed", fields)
		}
	}

	return &payload, nil
}
