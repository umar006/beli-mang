package domain

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
