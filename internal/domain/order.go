package domain

type Order struct {
	ID         string `json:"orderId" db:"id"`
	CreatedAt  int64  `json:"createdAt" db:"created_at"`
	MerchantID string `db:"merchant_id"`
}
