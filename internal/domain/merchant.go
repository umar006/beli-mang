package domain

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

type MerchantLocation struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type MerchantRequest struct {
	Name     string               `json:"name"`
	Category MerchantCategoryType `json:"merchantCategory"`
	ImageUrl string               `json:"imageUrl"`
	Location MerchantLocation     `json:"location"`
}

func (mr *MerchantRequest) NewMerchantFromDTO() Merchant {
	id, _ := gonanoid.New()
	createdAt := time.Now().UnixNano()

	return Merchant{
		ID:        id,
		CreatedAt: createdAt,
		Name:      mr.Name,
		Category:  mr.Category,
		ImageUrl:  mr.ImageUrl,
		Location: pgtype.Point{
			P: pgtype.Vec2{
				X: mr.Location.Latitude,
				Y: mr.Location.Longitude,
			},
			Valid: true,
		},
	}
}
