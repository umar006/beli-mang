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
	headers := ctx.GetReqHeaders()
	contentType := headers["Content-Type"]
	contentLength := headers["Content-Length"][0]

	if len(contentType) < 1 || contentLength == "0" {
		err := domain.NewErrBadRequest("empty file")
		return ctx.Status(err.Code).JSON(err)
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		err := domain.NewErrInternalServerError(err.Error())
		return ctx.Status(err.Code).JSON(err)
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
