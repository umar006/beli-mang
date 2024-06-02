package server

import (
	"beli-mang/internal/handler"
	"beli-mang/internal/middleware"
	"beli-mang/internal/repository"
	"beli-mang/internal/service"
	"beli-mang/internal/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var validate = validator.New()

func (s *FiberServer) RegisterFiberRoutes() {
	db := s.db.GetDB()

	s.App.Use(recover.New())

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	validate.RegisterValidation("url", validation.URL)

	userRepo := repository.NewUser()
	merchantRepo := repository.NewMerchantRepo()
	merchantItemRepo := repository.NewMerchantItemRepo()

	userService := service.NewUser(db, userRepo)
	merchantService := service.NewMerchantService(db, merchantRepo)
	merchantItemService := service.NewMerchantItemService(db, merchantRepo, merchantItemRepo)
	awsS3Service := service.NewAWSS3Service()

	userHandler := handler.NewUser(validate, userService)
	merchantHandler := handler.NewMerchantHandler(validate, merchantService)
	merchantItemHandler := handler.NewMerchantItemHandler(validate, merchantItemService)
	awsS3Handler := handler.NewAWSS3(awsS3Service)
	authMiddleware := middleware.NewAuth()

	awsS3 := s.App.Group("/image")
	awsS3.Use(authMiddleware.Auth())
	awsS3.Post("/", awsS3Handler.UploadImage)

	admin := s.App.Group("/admin")
	admin.Post("/register", userHandler.CreateAdmin)
	admin.Post("/login", userHandler.Login)

	users := s.App.Group("/users")
	users.Post("/register", userHandler.CreateCustomer)
	users.Post("/login", userHandler.Login)

	merchant := admin.Group("/merchants")
	merchant.Use(authMiddleware.Auth())
	merchant.Post("/", merchantHandler.CreateMerchant)
	merchant.Get("/", merchantHandler.GetMerchantList)
	merchant.Get("/nearby/:latlong", merchantHandler.GetMerchantListByLatLong)
	merchant.Post("/:merchantId/items", merchantItemHandler.CreateMerchantItem)
	merchant.Get("/:merchantId/items", merchantItemHandler.GetMerchantItemList)
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
