package service

import (
	"beli-mang/internal/domain"
	"context"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AWSS3Service interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, *fiber.Error)
}

type awsS3Service struct{}

func NewAWSS3Service() AWSS3Service {
	return &awsS3Service{}
}

var (
	awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	awsS3BucketName    = os.Getenv("AWS_S3_BUCKET_NAME")
	awsRegion          = os.Getenv("AWS_REGION")
)

func (a *awsS3Service) UploadImage(fileHeader *multipart.FileHeader) (string, *fiber.Error) {
	openFile, err := fileHeader.Open()
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}
	defer openFile.Close()

	fileExtension := filepath.Ext(fileHeader.Filename)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")))
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	client := s3.NewFromConfig(cfg)

	filename := uuid.New().String() + fileExtension
	client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(awsS3BucketName),
		Key:    aws.String(filename),
		Body:   openFile,
	})

	return filename, nil
}
