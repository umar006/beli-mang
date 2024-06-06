package domain

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type RoleType string

const (
	RoleAdmin    RoleType = "admin"
	RoleCustomer RoleType = "customer"
)

type User struct {
	ID        string   `json:"userId" db:"id"`
	CreatedAt int64    `json:"createdAt" db:"created_at"`
	Username  string   `json:"username" db:"username"`
	Email     string   `json:"email" db:"email"`
	Password  string   `json:"password" db:"password"`
	Role      RoleType `json:"role" db:"role"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type UserLocation struct {
	Latitude  float64 `json:"lat" validate:"required,number"`
	Longitude float64 `json:"long" validate:"required,number"`
}

type UserOrderItem struct {
	ItemID   string `json:"itemId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

type UserOrder struct {
	MerchantID      string          `json:"merchantId" validate:"required"`
	IsStartingPoint *bool           `json:"isStartingPoint" validate:"required"`
	OrderItems      []UserOrderItem `json:"items"`
}

type PriceEstimation struct {
	ID                             string `db:"id"`
	CreatedAt                      int64  `db:"created_at"`
	TotalPrice                     int    `db:"total_price"`
	EstimatedDeliveryTimeInMinutes int    `db:"delivery_time_in_minutes"`
}

type PriceEstimateRequest struct {
	UserLoc UserLocation `json:"userLocation"`
	Orders  []UserOrder  `json:"orders"`
}

type PriceEstimateResponse struct {
	ID                             string `json:"calculatedEstimatedId"`
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
}

func (ar *RegisterRequest) NewUserFromDTO() User {
	id, _ := gonanoid.New()
	createdAt := time.Now().UnixNano()

	return User{
		ID:        id,
		CreatedAt: createdAt,
		Username:  ar.Username,
		Email:     ar.Email,
		Password:  ar.Password,
	}
}
