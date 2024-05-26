package server

import (
	"beli-mang/internal/handler"
	"beli-mang/internal/repository"
	"beli-mang/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func (s *FiberServer) RegisterFiberRoutes() {
	db := s.db.GetDB()

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	userRepo := repository.NewUser()
	merchantRepo := repository.NewMerchantRepo()

	userService := service.NewUser(db, userRepo)
	merchantService := service.NewMerchantService(db, merchantRepo)

	userHandler := handler.NewUser(validate, userService)
	merchantHandler := handler.NewMerchantHandler(merchantService)

	admin := s.App.Group("/admin")
	admin.Post("/register", userHandler.CreateAdmin)
	admin.Post("/login", userHandler.Login)

	users := s.App.Group("/users")
	users.Post("/register", userHandler.CreateCustomer)
	users.Post("/login", userHandler.Login)

	merchant := admin.Group("/merchants")
	merchant.Post("/", merchantHandler.CreateMerchant)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
