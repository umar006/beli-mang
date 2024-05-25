package helper

import (
	"beli-mang/internal/domain"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func msgForTag(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	param := fe.Param()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("minimum %s is %s characters", field, param)
	case "max":
		return fmt.Sprintf("maximum %s is %s characters", field, param)
	case "email":
		return fmt.Sprintf("invalid %s format", field)
	}

	return "unhandled validation"
}

func ValidateRequest(err error) *fiber.Error {
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range err {
			return domain.NewErrBadRequest(msgForTag(fe))
		}
	}
	return domain.NewErrBadRequest(err.Error())
}
