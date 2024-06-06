package handler

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/helper"
	"beli-mang/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	CreateAdmin(ctx *fiber.Ctx) error
	CreateCustomer(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	GetPriceEstimation(ctx *fiber.Ctx) error
}

type userHandler struct {
	validator   *validator.Validate
	userService service.UserService
}

func NewUser(validator *validator.Validate, userService service.UserService) UserHandler {
	return &userHandler{
		validator:   validator,
		userService: userService,
	}
}

func (uh *userHandler) CreateAdmin(ctx *fiber.Ctx) error {
	var body domain.RegisterRequest
	ctx.BodyParser(&body)

	if err := uh.validator.Struct(&body); err != nil {
		err := helper.ValidateRequest(err)
		return ctx.Status(err.Code).JSON(err)
	}

	token, err := uh.userService.CreateAdmin(ctx.Context(), body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(201).JSON(map[string]string{"token": token})
}

func (uh *userHandler) CreateCustomer(ctx *fiber.Ctx) error {
	var body domain.RegisterRequest
	ctx.BodyParser(&body)

	if err := uh.validator.Struct(&body); err != nil {
		err := helper.ValidateRequest(err)
		return ctx.Status(err.Code).JSON(err)
	}

	token, err := uh.userService.CreateCustomer(ctx.Context(), body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(201).JSON(map[string]string{"token": token})
}

func (uh *userHandler) Login(ctx *fiber.Ctx) error {
	var body domain.LoginRequest
	ctx.BodyParser(&body)

	if err := uh.validator.Struct(&body); err != nil {
		err := helper.ValidateRequest(err)
		return ctx.Status(err.Code).JSON(err)
	}

	token, err := uh.userService.Login(ctx.Context(), body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(200).JSON(map[string]string{"token": token})
}

func (uh *userHandler) GetPriceEstimation(ctx *fiber.Ctx) error {
	var body domain.PriceEstimateRequest
	ctx.BodyParser(&body)

	_, err := uh.userService.GetPriceEstimation(ctx.Context(), body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return nil
}
