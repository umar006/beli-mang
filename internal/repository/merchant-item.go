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

type MerchantItemRepo interface {
	CreateMerchantItem(ctx context.Context, db *pgx.Conn, merchantId string, merchantItem domain.MerchantItem) error
	GetMerchantItemByID(ctx context.Context, db *pgx.Conn, itemID string) (domain.MerchantItemResponse, error)
	GetMerchantItemListByMerchantID(ctx context.Context, db *pgx.Conn, merchantId string, queryParams domain.MerchantItemQueryParams) ([]domain.MerchantItemResponse, *domain.Page, error)
	GetTotalMerchantItemListByMerchantID(ctx context.Context, db *pgx.Conn, merchantId string) (int, error)
	CheckItemExistsByID(ctx context.Context, db *pgx.Conn, itemID string) (bool, error)
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

func (mi *merchantItemRepo) GetMerchantItemListByMerchantID(ctx context.Context, db *pgx.Conn, merchantId string, queryParams domain.MerchantItemQueryParams) ([]domain.MerchantItemResponse, *domain.Page, error) {
	var whereParams []string
	var sortParams []string
	var limitOffsetParams []string
	var args []any
	args = append(args, merchantId)
	argPos := 2 // start from 2 because $1 for merchantId

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

	limit := 5
	parsedLimit, err := strconv.Atoi(queryParams.Limit)
	if queryParams.Limit != "" && err == nil {
		limit = parsedLimit
	}
	limitOffsetParams = append(limitOffsetParams, fmt.Sprintf("LIMIT $%d", argPos))
	args = append(args, limit)
	argPos++

	offset := 0
	parsedOffset, err := strconv.Atoi(queryParams.Offset)
	if queryParams.Limit != "" && err == nil {
		offset = parsedOffset
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

	query := `SELECT id, created_at, name, category, price, image_url
	 FROM merchant_items
	 WHERE merchant_id = $1`
	query += whereQuery
	query += sortQuery
	query += limitOffsetQuery

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}

	merchantItemList := []domain.MerchantItemResponse{}
	for rows.Next() {
		merchantItemFromDB := domain.MerchantItem{}
		err := rows.Scan(
			&merchantItemFromDB.ID, &merchantItemFromDB.CreatedAt, &merchantItemFromDB.Name,
			&merchantItemFromDB.Category, &merchantItemFromDB.Price, &merchantItemFromDB.ImageUrl,
		)
		if err != nil {
			return nil, nil, err
		}

		parsedCreatedAt := time.Unix(0, merchantItemFromDB.CreatedAt).Format(time.RFC3339)
		merchantItem := domain.MerchantItemResponse{
			ID:        merchantItemFromDB.ID,
			CreatedAt: parsedCreatedAt,
			Name:      merchantItemFromDB.Name,
			Category:  merchantItemFromDB.Category,
			Price:     merchantItemFromDB.Price,
			ImageUrl:  merchantItemFromDB.ImageUrl,
		}

		merchantItemList = append(merchantItemList, merchantItem)
	}

	totalMerchantItemList, err := mi.GetTotalMerchantItemListByMerchantID(ctx, db, merchantId)
	if err != nil {
		return nil, nil, err
	}

	page := &domain.Page{
		Limit:  limit,
		Offset: offset,
		Total:  totalMerchantItemList,
	}

	return merchantItemList, page, nil
}

func (mi *merchantItemRepo) GetTotalMerchantItemListByMerchantID(ctx context.Context, db *pgx.Conn, merchantId string) (int, error) {
	query := `SELECT count(*)
	FROM merchant_items
	WHERE merchant_id = $1`
	var total int
	err := db.QueryRow(ctx, query, merchantId).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (mi *merchantItemRepo) GetMerchantItemByID(ctx context.Context, db *pgx.Conn, itemID string) (domain.MerchantItemResponse, error) {
	query := `SELECT id, created_at, name, price, category, image_url
	FROM merchant_items
	WHERE id = $1`

	itemFromDB := domain.MerchantItem{}
	err := db.QueryRow(ctx, query, itemID).Scan(&itemFromDB.ID, &itemFromDB.CreatedAt,
		&itemFromDB.Name, &itemFromDB.Price, &itemFromDB.Category, &itemFromDB.ImageUrl,
	)
	if err != nil {
		return domain.MerchantItemResponse{}, err
	}

	parsedCreatedAt := time.Unix(0, itemFromDB.CreatedAt).Format(time.RFC3339)
	item := domain.MerchantItemResponse{
		ID:        itemFromDB.ID,
		CreatedAt: parsedCreatedAt,
		Name:      itemFromDB.Name,
		Category:  itemFromDB.Category,
		Price:     itemFromDB.Price,
		ImageUrl:  itemFromDB.ImageUrl,
	}

	return item, nil
}

func (mi *merchantItemRepo) CheckItemExistsByID(ctx context.Context, db *pgx.Conn, itemID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM merchant_items WHERE id = $1)`
	var exists bool
	err := db.QueryRow(ctx, query, itemID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
