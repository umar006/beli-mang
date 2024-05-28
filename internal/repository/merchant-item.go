package repository

import (
	"beli-mang/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MerchantItemRepo interface {
	CreateMerchantItem(ctx context.Context, db *pgx.Conn, merchantId string, merchantItem domain.MerchantItem) error
	GetMerchantItemListByMerchantID(ctx context.Context, db *pgx.Conn, merchantId string) ([]domain.MerchantItemResponse, error)
}

type merchantItemRepo struct{}

func NewMerchantItemRepo() MerchantItemRepo {
	return &merchantItemRepo{}
}

func (mi *merchantItemRepo) CreateMerchantItem(ctx context.Context, db *pgx.Conn, merchantId string, merchantItem domain.MerchantItem) error {
	query := `INSERT INTO merchant_items (id, created_at, name, category, price, image_url, merchant_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(
		ctx, query,
		merchantItem.ID, merchantItem.CreatedAt, merchantItem.Name,
		merchantItem.Category, merchantItem.Price, merchantItem.ImageUrl,
		merchantId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (mi *merchantItemRepo) GetMerchantItemListByMerchantID(ctx context.Context, db *pgx.Conn, merchantId string) ([]domain.MerchantItemResponse, error) {
	fmt.Println(merchantId)
	query := `SELECT id, created_at, name, category, price, image_url
	 FROM merchant_items
	 WHERE merchant_id = $1`
	rows, err := db.Query(ctx, query, merchantId)
	if err != nil {
		return nil, err
	}

	merchantItemList := []domain.MerchantItemResponse{}
	for rows.Next() {
		merchantItem := domain.MerchantItemResponse{}
		err := rows.Scan(
			&merchantItem.ID, &merchantItem.CreatedAt, &merchantItem.Name,
			&merchantItem.Category, &merchantItem.Price, &merchantItem.ImageUrl,
		)
		if err != nil {
			return nil, err
		}
		merchantItemList = append(merchantItemList, merchantItem)
	}

	return merchantItemList, nil
}
