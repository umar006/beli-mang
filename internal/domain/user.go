package domain

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type RoleType int8

const (
	Customer RoleType = iota
	Admin
)

func (r RoleType) String() string {
	switch r {
	case Customer:
		return "customer"
	case Admin:
		return "admin"
	}
	return "unknown role"
}

type User struct {
	ID        string   `json:"userId" db:"id"`
	CreatedAt int64    `json:"createdAt" db:"created_at"`
	Username  string   `json:"username" db:"username"`
	Email     string   `json:"email" db:"email"`
	Password  string   `json:"password" db:"password"`
	Role      RoleType `json:"role" db:"role"`
}

type AdminRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

func (ar *AdminRequest) NewUserFromDTO() User {
	id, _ := gonanoid.New()
	createdAt := time.Now().UnixNano()

	return User{
		ID:        id,
		CreatedAt: createdAt,
		Username:  ar.Username,
		Email:     ar.Email,
		Password:  ar.Password,
		Role:      Admin,
	}
}
