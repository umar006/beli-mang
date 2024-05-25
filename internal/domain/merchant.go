package domain

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MerchantCategoryType int8

const (
	SmallRestaurant MerchantCategoryType = iota
	MediumRestaurant
	LargeRestaurant
	MerchandiseRestaurant
	BoothKiosk
	ConvenienceStore
)

func (mc MerchantCategoryType) String() string {
	switch mc {
	case SmallRestaurant:
		return "SmallRestaurant"
	case MediumRestaurant:
		return "MediumRestaurant"
	case LargeRestaurant:
		return "LargeRestaurant"
	case MerchandiseRestaurant:
		return "MerchandiseRestaurant"
	case BoothKiosk:
		return "BoothKiosk"
	case ConvenienceStore:
		return "ConvenienceStore"
	}
	return "unknown merchant category"
}

type Merchant struct {
	ID        string               `json:"merchantId" db:"id"`
	CreatedAt int64                `json:"createdAt" db:"created_at"`
	Name      string               `json:"name" db:"name"`
	Category  MerchantCategoryType `json:"merchantCategory" db:"category"`
	ImageUrl  string               `json:"imageUrl" db:"image_url"`
	Location  pgtype.Point         `json:"location" db:"location"`
}
