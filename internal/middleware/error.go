package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog"

	domainerror "github.com/jaeyoung0509/titanbay-funds-api/internal/domain/error"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/response"
)

func NewErrorHandler(base zerolog.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		appErr := domainerror.As(err)
		if appErr == nil || appErr.Kind == domainerror.KindInternal {
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotFound {
				appErr = domainerror.NotFound("")
			}
		}
		if appErr == nil {
			appErr = domainerror.Internal(err)
		}

		if appErr.Kind == domainerror.KindInternal && appErr.Err != nil {
			base.Error().
				Str("request_id", requestid.FromContext(c)).
				Str("method", c.Method()).
				Str("path", c.OriginalURL()).
				Err(appErr.Err).
				Msg("internal error")
		}

		status := statusFor(appErr.Kind)
		return c.Status(status).JSON(response.NewErrorEnvelope(string(appErr.Kind), appErr.Message, appErr.Fields))
	}
}

func statusFor(kind domainerror.Kind) int {
	switch kind {
	case domainerror.KindValidation:
		return fiber.StatusBadRequest
	case domainerror.KindNotFound:
		return fiber.StatusNotFound
	case domainerror.KindConflict:
		return fiber.StatusConflict
	case domainerror.KindInternal:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}
