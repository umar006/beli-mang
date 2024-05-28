package domain

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type ItemCategoryType string

const (
	ItemBeverage   ItemCategoryType = "Beverage"
	ItemFood       ItemCategoryType = "Food"
	ItemSnack      ItemCategoryType = "Snack"
	ItemCondiments ItemCategoryType = "Condiments"
	ItemAdditions  ItemCategoryType = "Additions"
)

type MerchantItem struct {
	ID        string           `json:"itemId" db:"id"`
	CreatedAt int64            `json:"createdAt" db:"created_at"`
	Name      string           `json:"name" db:"name"`
	Category  ItemCategoryType `json:"productCategory" db:"category"`
	Price     int64            `json:"price" db:"price"`
	ImageUrl  string           `json:"imageUrl" db:"image_url"`
}

type MerchantItemRequest struct {
	Name     string           `json:"name" validate:"required,min=2,max=30"`
	Category ItemCategoryType `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	Price    int64            `json:"price" validate:"required,min=1"`
	ImageUrl string           `json:"imageUrl" validate:"required,url"`
}

func (mi *MerchantItemRequest) NewMerchantItemFromDTO() MerchantItem {
	id, _ := gonanoid.New()
	createdAt := time.Now().UnixNano()

	return MerchantItem{
		ID:        id,
		CreatedAt: createdAt,
		Name:      mi.Name,
		Category:  mi.Category,
		Price:     mi.Price,
		ImageUrl:  mi.ImageUrl,
	}
}
