package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(response.NewHealth())
}
