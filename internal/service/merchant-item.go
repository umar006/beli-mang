package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type MerchantItemService interface {
	CreateMerchantItem(ctx context.Context, merchantId string, body domain.MerchantItemRequest) (string, *fiber.Error)
}

type merchantItemService struct {
	db               *pgx.Conn
	merchantItemRepo repository.MerchantItemRepo
}

func NewMerchantItemService(db *pgx.Conn, merchantItemRepo repository.MerchantItemRepo) MerchantItemService {
	return &merchantItemService{
		db:               db,
		merchantItemRepo: merchantItemRepo,
	}
}

func (mi *merchantItemService) CreateMerchantItem(ctx context.Context, merchantId string, body domain.MerchantItemRequest) (string, *fiber.Error) {
	merchantItem := body.NewMerchantItemFromDTO()

	err := mi.merchantItemRepo.CreateMerchantItem(ctx, mi.db, merchantId, merchantItem)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	return merchantItem.ID, nil
}
