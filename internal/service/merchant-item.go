package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/repository"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MerchantItemService interface {
	CreateMerchantItem(ctx context.Context, merchantId string, body domain.MerchantItemRequest) (string, *fiber.Error)
	GetMerchantItemListByMerchantID(ctx context.Context, merchantId string) ([]domain.MerchantItemResponse, *fiber.Error)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return "", domain.NewErrNotFound("merchant is not found")
			}
		}
		return "", domain.NewErrInternalServerError(err.Error())
	}

	return merchantItem.ID, nil
}

func (mi *merchantItemService) GetMerchantItemListByMerchantID(ctx context.Context, merchantId string) ([]domain.MerchantItemResponse, *fiber.Error) {
	merchantItemList, err := mi.merchantItemRepo.GetMerchantItemListByMerchantID(ctx, mi.db, merchantId)
	if err != nil {
		return nil, domain.NewErrInternalServerError(err.Error())
	}

	return merchantItemList, nil
}
