package handler

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/helper"
	"beli-mang/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MerchantItemHandler interface {
	CreateMerchantItem(ctx *fiber.Ctx) error
}

type merchantItemHandler struct {
	validator           *validator.Validate
	merchantItemService service.MerchantItemService
}

func NewMerchantItemHandler(validator *validator.Validate, merchantItemService service.MerchantItemService) MerchantItemHandler {
	return &merchantItemHandler{
		validator:           validator,
		merchantItemService: merchantItemService,
	}
}

func (mi *merchantItemHandler) CreateMerchantItem(ctx *fiber.Ctx) error {
	merchantId := ctx.Params("merchantId", "")

	var body domain.MerchantItemRequest
	ctx.BodyParser(&body)

	if err := mi.validator.Struct(body); err != nil {
		err := helper.ValidateRequest(err)
		return ctx.Status(err.Code).JSON(err)
	}

	merchantItemId, err := mi.merchantItemService.CreateMerchantItem(ctx.Context(), merchantId, body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(201).JSON(map[string]string{"itemId": merchantItemId})
}
