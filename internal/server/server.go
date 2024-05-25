package server

import (
	"github.com/gofiber/fiber/v2"

	"beli-mang/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "beli-mang",
			AppName:      "beli-mang",
		}),

		db: database.New(),
	}

	return server
}
