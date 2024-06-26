package domain

import "github.com/gofiber/fiber/v2"

func NewErrInternalServerError(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}

func NewErrConflict(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusConflict, msg)
}

func NewErrBadRequest(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}

func NewErrNotFound(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusNotFound, msg)
}

func NewErrUnauthorized(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusUnauthorized, msg)
}
