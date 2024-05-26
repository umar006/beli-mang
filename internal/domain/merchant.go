package domain

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MerchantCategoryType string

const (
	MerchantSmallRestaurant       MerchantCategoryType = "SmallRestaurant"
	MerchantMediumRestaurant      MerchantCategoryType = "MediumRestaurant"
	MerchantLargeRestaurant       MerchantCategoryType = "LargeRestaurant"
	MerchantMerchandiseRestaurant MerchantCategoryType = "MerchandiseRestaurant"
	MerchantBoothKiosk            MerchantCategoryType = "BoothKiosk"
	MerchantConvenienceStore      MerchantCategoryType = "ConvenienceStore"
)

type Merchant struct {
	ID        string               `json:"merchantId" db:"id"`
	CreatedAt int64                `json:"createdAt" db:"created_at"`
	Name      string               `json:"name" db:"name"`
	Category  MerchantCategoryType `json:"merchantCategory" db:"category"`
	ImageUrl  string               `json:"imageUrl" db:"image_url"`
	Location  pgtype.Point         `json:"location" db:"location"`
}
