package handler

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/helper"
	"beli-mang/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type MerchantHandler interface {
	CreateMerchant(ctx *fiber.Ctx) error
}

type merchantHandler struct {
	validator       *validator.Validate
	merchantService service.MerchantService
}

func NewMerchantHandler(validator *validator.Validate, merchantService service.MerchantService) MerchantHandler {
	return &merchantHandler{
		validator:       validator,
		merchantService: merchantService,
	}
}

func (mh *merchantHandler) CreateMerchant(ctx *fiber.Ctx) error {
	var body domain.MerchantRequest
	ctx.BodyParser(&body)

	if err := mh.validator.Struct(&body); err != nil {
		err := helper.ValidateRequest(err)
		return ctx.Status(err.Code).JSON(err)
	}

	merchantID, err := mh.merchantService.CreateMerchant(ctx.Context(), body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(201).JSON(map[string]string{"merchantId": merchantID})
}
