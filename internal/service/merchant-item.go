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
	GetMerchantItemListByMerchantID(ctx context.Context, merchantId string, queryParams domain.MerchantItemQueryParams) ([]domain.MerchantItemResponse, *domain.Page, *fiber.Error)
}

type merchantItemService struct {
	db               *pgx.Conn
	merchantRepo     repository.MerchantRepo
	merchantItemRepo repository.MerchantItemRepo
}

func NewMerchantItemService(db *pgx.Conn, merchantRepo repository.MerchantRepo, merchantItemRepo repository.MerchantItemRepo) MerchantItemService {
	return &merchantItemService{
		db:               db,
		merchantRepo:     merchantRepo,
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

func (mi *merchantItemService) GetMerchantItemListByMerchantID(ctx context.Context, merchantId string, queryParams domain.MerchantItemQueryParams) ([]domain.MerchantItemResponse, *domain.Page, *fiber.Error) {
	exists, err := mi.merchantRepo.CheckMerchantExistsByMerchantID(ctx, mi.db, merchantId)
	if err != nil {
		return nil, nil, domain.NewErrInternalServerError(err.Error())
	}
	if !exists {
		return nil, nil, domain.NewErrNotFound("merchant is not found")
	}

	merchantItemList, page, err := mi.merchantItemRepo.GetMerchantItemListByMerchantID(ctx, mi.db, merchantId, queryParams)
	if err != nil {
		return nil, nil, domain.NewErrInternalServerError(err.Error())
	}

	return merchantItemList, page, nil
}
