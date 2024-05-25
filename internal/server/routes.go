package server

import (
	"beli-mang/internal/handler"
	"beli-mang/internal/repository"
	"beli-mang/internal/service"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	db := s.db.GetDB()

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	userRepo := repository.NewUser()

	userService := service.NewUser(db, userRepo)

	userHandler := handler.NewUser(userService)

	admin := s.App.Group("/admin")
	admin.Post("/register", userHandler.CreateAdmin)
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
