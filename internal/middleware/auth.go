package middleware

import (
	"beli-mang/internal/domain"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type Auth interface {
	Auth() fiber.Handler
}

type auth struct{}

func NewAuth() Auth {
	return &auth{}
}

var jwtSecret = os.Getenv("JWT_SECRET")

func (a *auth) Auth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(jwtSecret)},
		ErrorHandler: customAuthError,
	})
}

func customAuthError(ctx *fiber.Ctx, err error) error {
	if err.Error() == jwtware.ErrJWTMissingOrMalformed.Error() {
		missingOrMalformed := domain.NewErrUnauthorized(jwtware.ErrJWTMissingOrMalformed.Error())
		return ctx.Status(missingOrMalformed.Code).JSON(missingOrMalformed)
	}
	invalidOrExpired := domain.NewErrUnauthorized("Invalid or expired JWT")
	return ctx.Status(invalidOrExpired.Code).JSON(invalidOrExpired)
}
