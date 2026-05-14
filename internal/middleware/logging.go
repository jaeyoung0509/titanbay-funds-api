package middleware

import (
	zerologmw "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog"
)

func NewRequestLogger(base zerolog.Logger) fiber.Handler {
	logger := base

	return zerologmw.New(zerologmw.Config{
		Logger: &logger,
		GetLogger: func(c fiber.Ctx) zerolog.Logger {
			return base.With().
				Str("request_id", requestid.FromContext(c)).
				Logger()
		},
		Fields: []string{
			"latency",
			"status",
			"method",
			"url",
			"error",
		},
		FieldsSnakeCase: true,
	})
}
