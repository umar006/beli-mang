package repository

import (
	"beli-mang/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, db *pgx.Conn, merchant domain.Merchant) error
}

type merchantRepo struct{}

func NewMerchantRepo() MerchantRepo {
	return &merchantRepo{}
}

func (mr *merchantRepo) CreateMerchant(ctx context.Context, db *pgx.Conn, merchant domain.Merchant) error {
	query := `INSERT INTO merchants (id, created_at, name, category, image_url, location)
				VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(
		ctx, query,
		merchant.ID, merchant.CreatedAt, merchant.Name,
		merchant.Category, merchant.ImageUrl, merchant.Location,
	)
	if err != nil {
		return err
	}
	return nil
}
