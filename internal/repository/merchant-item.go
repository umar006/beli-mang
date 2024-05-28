package repository

import (
	"beli-mang/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type MerchantItemRepo interface {
	CreateMerchantItem(ctx context.Context, db *pgx.Conn, merchantId string, merchantItem domain.MerchantItem) error
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
