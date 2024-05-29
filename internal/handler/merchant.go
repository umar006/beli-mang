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
	GetMerchantList(ctx *fiber.Ctx) error
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

func (mh *merchantHandler) GetMerchantList(ctx *fiber.Ctx) error {
	var queryParams domain.MerchantQueryParams
	ctx.QueryParser(&queryParams)

	merchantList, page, err := mh.merchantService.GetMerchantList(ctx.Context(), queryParams)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	response := domain.SuccessResponse{
		Data: merchantList,
		Meta: page,
	}

	return ctx.Status(200).JSON(response)
}
