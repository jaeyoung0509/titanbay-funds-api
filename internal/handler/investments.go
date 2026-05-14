package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/request"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/response"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/usecase"
)

type InvestmentHandler struct {
	service usecase.InvestmentService
}

func NewInvestmentHandler(service usecase.InvestmentService) *InvestmentHandler {
	return &InvestmentHandler{service: service}
}

func (h *InvestmentHandler) ListInvestments(c fiber.Ctx) error {
	fundID, err := parseIDField("fund_id", c.Params("fund_id"))
	if err != nil {
		return err
	}

	items, err := h.service.ListInvestmentsByFund(c, fundID)
	if err != nil {
		return err
	}

	out := make([]response.Investment, 0, len(items))
	for _, item := range items {
		out = append(out, response.NewInvestment(item))
	}
	return c.JSON(out)
}

func (h *InvestmentHandler) CreateInvestment(c fiber.Ctx) error {
	fundID, err := parseIDField("fund_id", c.Params("fund_id"))
	if err != nil {
		return err
	}

	req, err := request.BindAndValidate[request.CreateInvestmentRequest](c)
	if err != nil {
		return err
	}

	input, err := req.ToInput(fundID)
	if err != nil {
		return err
	}

	item, err := h.service.CreateInvestment(c, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewInvestment(item))
}
