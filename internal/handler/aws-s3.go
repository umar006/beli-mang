package handler

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AWSS3Handler interface {
	UploadImage(ctx *fiber.Ctx) error
}

type awsS3 struct {
	awsS3Service service.AWSS3Service
}

func NewAWSS3(awsS3Service service.AWSS3Service) AWSS3Handler {
	return &awsS3{
		awsS3Service: awsS3Service,
	}
}

func (aws *awsS3) UploadImage(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(400).JSON(err)
	}

	var filename string
	if f, err := aws.awsS3Service.UploadImage(file); err != nil {
		return ctx.Status(400).JSON(err)
	} else {
		filename = f
	}

	response := domain.SuccessResponse{
		Data: map[string]string{"imageUrl": filename},
	}

	return ctx.Status(200).JSON(response)
}
