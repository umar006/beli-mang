package domain

type ItemCategoryType int8

const (
	Beverage ItemCategoryType = iota
	Food
	Snack
	Condiments
	Additions
)

func (ic ItemCategoryType) String() string {
	switch ic {
	case Beverage:
		return "Beverage"
	case Food:
		return "Food"
	case Snack:
		return "Snack"
	case Condiments:
		return "Condiments"
	case Additions:
		return "Additions"
	}
	return "unknown item category"
}

type MerchantItem struct {
	ID        string           `json:"itemId" db:"id"`
	CreatedAt int64            `json:"createdAt" db:"created_at"`
	Name      string           `json:"name" db:"name"`
	Category  ItemCategoryType `json:"productCategory" db:"category"`
	Price     int64            `json:"price" db:"price"`
	ImageUrl  string           `json:"imageUrl" db:"image_url"`
}
