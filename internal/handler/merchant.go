package handler

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MerchantHandler interface {
	CreateMerchant(ctx *fiber.Ctx) error
}

type merchantHandler struct {
	merchantService service.MerchantService
}

func NewMerchantHandler(merchantService service.MerchantService) MerchantHandler {
	return &merchantHandler{
		merchantService: merchantService,
	}
}

func (mh *merchantHandler) CreateMerchant(ctx *fiber.Ctx) error {
	var body domain.MerchantRequest
	ctx.BodyParser(&body)

	merchantID, err := mh.merchantService.CreateMerchant(ctx.Context(), body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(201).JSON(map[string]string{"merchantId": merchantID})
}
