package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/repository"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type MerchantService interface {
	CreateMerchant(ctx context.Context, body domain.MerchantRequest) (string, *fiber.Error)
	GetMerchantList(ctx context.Context, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, *fiber.Error)
	GetMerchantListByLatLong(ctx context.Context, latlong string, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, *fiber.Error)
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

func (ms *merchantService) GetMerchantList(ctx context.Context, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, *fiber.Error) {
	merchantList, page, err := ms.merchantRepo.GetMerchantList(ctx, ms.db, queryParams)
	if err != nil {
		return nil, nil, domain.NewErrInternalServerError(err.Error())
	}

	return merchantList, page, nil
}

func (ms *merchantService) GetMerchantListByLatLong(ctx context.Context, latlong string, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, *fiber.Error) {
	latlongAsSlices := strings.Split(latlong, ",")
	merchantList, page, err := ms.merchantRepo.GetMerchantListByLatLong(ctx, ms.db, latlongAsSlices, queryParams)
	if err != nil {
		return nil, nil, domain.NewErrInternalServerError(err.Error())
	}

	return merchantList, page, nil
}
