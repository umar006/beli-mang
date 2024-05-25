package handler

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	CreateAdmin(ctx *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

func NewUser(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (uh *userHandler) CreateAdmin(ctx *fiber.Ctx) error {
	var body domain.AdminRequest
	ctx.BodyParser(&body)

	token, err := uh.userService.CreateAdmin(ctx.Context(), body)
	if err != nil {
		return ctx.Status(err.Code).JSON(err)
	}

	return ctx.Status(201).JSON(map[string]string{"token": token})
}
