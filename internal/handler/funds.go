package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/domain/vo"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/request"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/response"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/usecase"
)

type FundHandler struct {
	service usecase.FundService
}

func NewFundHandler(service usecase.FundService) *FundHandler {
	return &FundHandler{service: service}
}

func (h *FundHandler) ListFunds(c fiber.Ctx) error {
	items, err := h.service.ListFunds(c)
	if err != nil {
		return err
	}

	out := make([]response.Fund, 0, len(items))
	for _, item := range items {
		out = append(out, response.NewFund(item))
	}
	return c.JSON(out)
}

func (h *FundHandler) CreateFund(c fiber.Ctx) error {
	req, err := request.BindAndValidate[request.CreateFundRequest](c)
	if err != nil {
		return err
	}

	input, err := req.ToInput()
	if err != nil {
		return err
	}

	item, err := h.service.CreateFund(c, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewFund(item))
}

func (h *FundHandler) UpdateFund(c fiber.Ctx) error {
	req, err := request.BindAndValidate[request.UpdateFundRequest](c)
	if err != nil {
		return err
	}

	input, err := req.ToInput()
	if err != nil {
		return err
	}

	item, err := h.service.UpdateFund(c, input)
	if err != nil {
		return err
	}

	return c.JSON(response.NewFund(item))
}

func (h *FundHandler) GetFund(c fiber.Ctx) error {
	id, err := parseIDField("id", c.Params("id"))
	if err != nil {
		return err
	}

	item, err := h.service.GetFund(c, id)
	if err != nil {
		return err
	}

	return c.JSON(response.NewFund(item))
}

func parseIDField(field, raw string) (vo.ID, error) {
	id, err := vo.ParseID(raw)
	if err != nil {
		return vo.ID{}, validationError(field)
	}
	return id, nil
}
