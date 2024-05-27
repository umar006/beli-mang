package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type MerchantService interface {
	CreateMerchant(ctx context.Context, body domain.MerchantRequest) (string, *fiber.Error)
	GetMerchantList(ctx context.Context) ([]domain.MerchantResponse, *fiber.Error)
}

type merchantService struct {
	db           *pgx.Conn
	merchantRepo repository.MerchantRepo
}

func NewMerchantService(db *pgx.Conn, merchantRepo repository.MerchantRepo) MerchantService {
	return &merchantService{
		db:           db,
		merchantRepo: merchantRepo,
	}
}

func (ms *merchantService) CreateMerchant(ctx context.Context, body domain.MerchantRequest) (string, *fiber.Error) {
	merchant := body.NewMerchantFromDTO()
	err := ms.merchantRepo.CreateMerchant(ctx, ms.db, merchant)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	return merchant.ID, nil
}

func (ms *merchantService) GetMerchantList(ctx context.Context) ([]domain.MerchantResponse, *fiber.Error) {
	merchantList, err := ms.merchantRepo.GetMerchantList(ctx, ms.db)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	return merchantList, nil
}
