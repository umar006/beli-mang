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
	GetMerchantByID(ctx context.Context, db *pgx.Conn, merchantID string) (domain.MerchantResponse, error)
	GetMerchantList(ctx context.Context, db *pgx.Conn, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, error)
	GetMerchantListByIDs(ctx context.Context, db *pgx.Conn, merchantIDs []string, userLoc domain.UserLocation) ([]domain.MerchantResponse, error)
	GetMerchantListByLatLong(ctx context.Context, db *pgx.Conn, latlong []string, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, error)
	GetTotalMerchantList(ctx context.Context, db *pgx.Conn) (int, error)
	CheckMerchantExistsByMerchantID(ctx context.Context, db *pgx.Conn, merchantID string) (bool, error)
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

func (mr *merchantRepo) GetMerchantList(ctx context.Context, db *pgx.Conn, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, error) {
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
	if queryParams.Offset != "" && err == nil {
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

	query := `SELECT id, created_at, name, category, image_url, location
	FROM merchants`
	query += whereQuery
	query += sortQuery
	query += limitOffsetQuery

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
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
	totalMerchantList, err := mr.GetTotalMerchantList(ctx, db)
	if err != nil {
		return nil, nil, err
	}

	page := domain.Page{
		Limit:  limit,
		Offset: offset,
		Total:  totalMerchantList,
	}

	return merchantList, &page, nil
}

func (mr *merchantRepo) GetTotalMerchantList(ctx context.Context, db *pgx.Conn) (int, error) {
	query := `SELECT count(*) FROM merchants`
	var total int
	err := db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (mr *merchantRepo) CheckMerchantExistsByMerchantID(ctx context.Context, db *pgx.Conn, merchantID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM merchants WHERE id = $1)`
	var exists bool
	err := db.QueryRow(ctx, query, merchantID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (mr *merchantRepo) GetMerchantListByLatLong(ctx context.Context, db *pgx.Conn, latlong []string, queryParams domain.MerchantQueryParams) ([]domain.MerchantResponse, *domain.Page, error) {
	var whereParams []string
	var limitOffsetParams []string
	args := []any{latlong[1], latlong[0]}
	argPos := 3 // start from 3 because 1 2 taken by latlong

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
	if queryParams.Offset != "" && err == nil {
		offset = parsedOffset
	}
	limitOffsetParams = append(limitOffsetParams, fmt.Sprintf("OFFSET $%d", argPos))
	args = append(args, offset)
	argPos++

	var whereQuery string
	if len(whereParams) > 0 {
		whereQuery = "\nWHERE " + strings.Join(whereParams, " AND ")
	}
	var limitOffsetQuery string
	limitOffsetQuery = "\n" + strings.Join(limitOffsetParams, " ")

	subqueryMerchant := fmt.Sprintf(`
		WITH merchantSq AS (
			SELECT m.id, m.created_at, m.name, m.category, m.image_url, m.location
			FROM merchants m
			%s
			%s
		)
	`, whereQuery, limitOffsetQuery)
	query := `SELECT m.id, m.created_at, m.name, m.category, m.image_url, m.location,
					mi.id, mi.created_at, mi.name, mi.category, mi.price, mi.image_url
	FROM merchantSq m
	JOIN merchant_items mi ON mi.merchant_id = m.id
	ORDER BY (m.location <@> point($1,$2)) ASC`
	query = subqueryMerchant + query

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}

	mapMerchantList := make(map[string]domain.MerchantResponse)
	mapItemList := make(map[string][]domain.MerchantItemResponse)
	for rows.Next() {
		merchantFromDB := domain.Merchant{}
		itemFromDB := domain.MerchantItem{}
		rows.Scan(&merchantFromDB.ID, &merchantFromDB.CreatedAt, &merchantFromDB.Name,
			&merchantFromDB.Category, &merchantFromDB.ImageUrl, &merchantFromDB.Location,
			&itemFromDB.ID, &itemFromDB.CreatedAt, &itemFromDB.Name, &itemFromDB.Category,
			&itemFromDB.Price, &itemFromDB.ImageUrl,
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

		parsedCreatedAt = time.Unix(0, itemFromDB.CreatedAt).Format(time.RFC3339)
		item := domain.MerchantItemResponse{
			ID:        itemFromDB.ID,
			CreatedAt: parsedCreatedAt,
			Name:      itemFromDB.Name,
			Category:  itemFromDB.Category,
			Price:     itemFromDB.Price,
			ImageUrl:  itemFromDB.ImageUrl,
		}

		mapMerchantList[merchant.ID] = merchant
		mapItemList[merchant.ID] = append(mapItemList[merchant.ID], item)
	}

	merchantList := []domain.MerchantResponse{}
	for k, v := range mapMerchantList {
		items := mapItemList[k]
		v.MerchantItems = items
		merchantList = append(merchantList, v)
	}

	return merchantList, nil, nil
}

func (mr *merchantRepo) GetMerchantByID(ctx context.Context, db *pgx.Conn, merchantID string) (domain.MerchantResponse, error) {
	query := `SELECT id, created_at, name, category, image_url, location
	FROM merchants
	WHERE id = $1`
	merchantFromDB := domain.Merchant{}
	err := db.QueryRow(ctx, query, merchantID).Scan(&merchantFromDB.ID, &merchantFromDB.CreatedAt,
		&merchantFromDB.Name, &merchantFromDB.Category, &merchantFromDB.ImageUrl, &merchantFromDB.Location,
	)
	if err != nil {
		return domain.MerchantResponse{}, err
	}

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

	return merchant, nil
}

func (mr *merchantRepo) GetMerchantListByIDs(ctx context.Context, db *pgx.Conn, merchantIDs []string, userLoc domain.UserLocation) ([]domain.MerchantResponse, error) {
	query := `SELECT m.id, m.created_at, m.name, m.category, m.image_url, m.location, (m.location <@> point($2, $3)) * 1.60934 AS distance
			FROM merchants m
			WHERE id = ANY($1)`
	rows, err := db.Query(ctx, query, merchantIDs, userLoc.Longitude, userLoc.Latitude)
	if err != nil {
		return nil, err
	}

	merchantList := []domain.MerchantResponse{}
	for rows.Next() {
		merchantFromDB := domain.Merchant{}
		var distance float64
		rows.Scan(&merchantFromDB.ID, &merchantFromDB.CreatedAt, &merchantFromDB.Name,
			&merchantFromDB.Category, &merchantFromDB.ImageUrl, &merchantFromDB.Location,
			&distance,
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
			DistanceFromUser: distance,
		}
		merchantList = append(merchantList, merchant)
	}

	return merchantList, nil
}
