package repository

import (
	"beli-mang/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, db *pgx.Conn, merchant domain.Merchant) error
	GetMerchantList(ctx context.Context, db *pgx.Conn) ([]domain.MerchantResponse, error)
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

func (mr *merchantRepo) GetMerchantList(ctx context.Context, db *pgx.Conn) ([]domain.MerchantResponse, error) {
	query := `SELECT id, created_at, name, category, image_url, location
	FROM merchants`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	merchantList := []domain.MerchantResponse{}
	for rows.Next() {
		merchantFromDB := domain.Merchant{}
		rows.Scan(&merchantFromDB.ID, &merchantFromDB.CreatedAt, &merchantFromDB.Name,
			&merchantFromDB.Category, &merchantFromDB.ImageUrl, &merchantFromDB.Location,
		)

		merchant := domain.MerchantResponse{
			ID:        merchantFromDB.ID,
			CreatedAt: merchantFromDB.CreatedAt,
			Name:      merchantFromDB.Name,
			Category:  merchantFromDB.Category,
			ImageUrl:  merchantFromDB.ImageUrl,
			Location: domain.MerchantLocation{
				Latitude:  merchantFromDB.Location.P.X,
				Longitude: merchantFromDB.Location.P.Y,
			},
		}
		merchantList = append(merchantList, merchant)
	}

	return merchantList, nil
}
