package repository

import (
	"beli-mang/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	CreateAdmin(ctx context.Context, db *pgx.Conn, admin domain.User) error
}

type userRepo struct{}

func NewUser() UserRepo {
	return &userRepo{}
}

func (ur *userRepo) CreateAdmin(ctx context.Context, db *pgx.Conn, admin domain.User) error {
	query := `INSERT INTO users (id, created_at, username, email, password, role)
                VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(ctx, query, admin.ID, admin.CreatedAt, admin.Username,
		admin.Email, admin.Password, admin.Role)
	if err != nil {
		return err
	}
	return nil
}
