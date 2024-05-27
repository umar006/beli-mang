package repository

import (
	"beli-mang/internal/domain"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, db *pgx.Conn, merchant domain.Merchant) error
	GetMerchantList(ctx context.Context, db *pgx.Conn, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, error)
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

func (mr *merchantRepo) GetMerchantList(ctx context.Context, db *pgx.Conn, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, error) {
	var whereParams []string
	var sortParams []string
	var limitOffsetParams []string
	var args []any
	argPos := 1

	if queryParams.ID != "" {
		whereParams = append(whereParams, fmt.Sprintf("id = $%d", argPos))
		args = append(args, queryParams.ID)
		argPos++
	}

	if queryParams.Name != "" {
		whereParams = append(whereParams, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+queryParams.Name+"%")
		argPos++
	}

	if queryParams.Category != "" {
		whereParams = append(whereParams, fmt.Sprintf("category = $%d", argPos))
		args = append(args, queryParams.Category)
		argPos++
	}

	if queryParams.CreatedAt == "asc" || queryParams.CreatedAt == "desc" {
		sortParams = append(sortParams, fmt.Sprintf("created_at %s", queryParams.CreatedAt))
	}

	limit := "5"
	_, err := strconv.Atoi(queryParams.Limit)
	if queryParams.Limit != "" && err == nil {
		limit = queryParams.Limit
	}
	limitOffsetParams = append(limitOffsetParams, fmt.Sprintf("LIMIT $%d", argPos))
	args = append(args, limit)
	argPos++

	offset := "0"
	_, err = strconv.Atoi(queryParams.Offset)
	if queryParams.Limit != "" && err == nil {
		offset = queryParams.Offset
	}
	limitOffsetParams = append(limitOffsetParams, fmt.Sprintf("OFFSET $%d", argPos))
	args = append(args, offset)
	argPos++

	var whereQuery string
	if len(whereParams) > 0 {
		whereQuery = "\nWHERE " + strings.Join(whereParams, " AND ")
	}
	var sortQuery string
	if len(sortParams) > 0 {
		sortQuery = "\nORDER BY " + strings.Join(sortParams, ", ")
	}
	var limitOffsetQuery string
	limitOffsetQuery = "\n" + strings.Join(limitOffsetParams, " ")

	query := `SELECT id, created_at, name, category, image_url, location
	FROM merchants`
	query += whereQuery
	query += sortQuery
	query += limitOffsetQuery

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	merchantList := []domain.MerchantResponse{}
	for rows.Next() {
		merchantFromDB := domain.Merchant{}
		rows.Scan(&merchantFromDB.ID, &merchantFromDB.CreatedAt, &merchantFromDB.Name,
			&merchantFromDB.Category, &merchantFromDB.ImageUrl, &merchantFromDB.Location,
		)

		parsedCreatedAt := time.Unix(0, merchantFromDB.CreatedAt).Format(time.RFC3339)
		merchant := domain.MerchantResponse{
			ID:        merchantFromDB.ID,
			CreatedAt: parsedCreatedAt,
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
