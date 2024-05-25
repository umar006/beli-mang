package domain

type OrderItem struct {
	OrderID        string `db:"order_id"`
	MerchantItemID string `db:"merchant_item_id"`
}
