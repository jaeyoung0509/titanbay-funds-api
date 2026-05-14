package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/request"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/response"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/usecase"
)

type InvestorHandler struct {
	service usecase.InvestorService
}

func NewInvestorHandler(service usecase.InvestorService) *InvestorHandler {
	return &InvestorHandler{service: service}
}

func (h *InvestorHandler) ListInvestors(c fiber.Ctx) error {
	items, err := h.service.ListInvestors(c)
	if err != nil {
		return err
	}

	out := make([]response.Investor, 0, len(items))
	for _, item := range items {
		out = append(out, response.NewInvestor(item))
	}
	return c.JSON(out)
}

func (h *InvestorHandler) CreateInvestor(c fiber.Ctx) error {
	req, err := request.BindAndValidate[request.CreateInvestorRequest](c)
	if err != nil {
		return err
	}

	input, err := req.ToInput()
	if err != nil {
		return err
	}

	item, err := h.service.CreateInvestor(c, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewInvestor(item))
}
