package handler

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MerchantItemHandler interface {
	CreateMerchantItem(ctx *fiber.Ctx) error
}

type merchantItemHandler struct {
	merchantItemService service.MerchantItemService
}

func NewMerchantItemHandler(merchantItemService service.MerchantItemService) MerchantItemHandler {
	return &merchantItemHandler{
		merchantItemService: merchantItemService,
	}
}

func (mi *merchantItemHandler) CreateMerchantItem(ctx *fiber.Ctx) error {
	merchantId := ctx.Params("merchantId", "")

	var body domain.MerchantItemRequest
	ctx.BodyParser(&body)

	merchantItemId, err := mi.merchantItemService.CreateMerchantItem(ctx.Context(), merchantId, body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(201).JSON(map[string]string{"itemId": merchantItemId})
}
