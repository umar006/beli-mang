package domain

import "github.com/gofiber/fiber/v2"

func NewErrInternalServerError(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}
