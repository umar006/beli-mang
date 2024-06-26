package repository

import (
	"beli-mang/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	Create(ctx context.Context, db *pgx.Conn, user domain.User) error
	CreatePriceEstimate(ctx context.Context, db *pgx.Conn, priceEstimate domain.PriceEstimation) error
	GetUserByUsername(ctx context.Context, db *pgx.Conn, username string) (domain.User, error)
}

type userRepo struct{}

func NewUser() UserRepo {
	return &userRepo{}
}

func (ur *userRepo) Create(ctx context.Context, db *pgx.Conn, user domain.User) error {
	query := `INSERT INTO users (id, created_at, username, email, password, role)
                VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(ctx, query, user.ID, user.CreatedAt, user.Username,
		user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepo) GetUserByUsername(ctx context.Context, db *pgx.Conn, username string) (domain.User, error) {
	query := `SELECT id, created_at, username, email, password, role
				FROM users
				WHERE username = $1`
	var user domain.User
	err := db.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.CreatedAt, &user.Username,
			&user.Email, &user.Password, &user.Role)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ur *userRepo) CreatePriceEstimate(ctx context.Context, db *pgx.Conn, priceEstimate domain.PriceEstimation) error {
	query := `INSERT INTO price_estimations (id, created_at, total_price, delivery_time_in_minutes)
				VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(ctx, query, priceEstimate.ID, priceEstimate.CreatedAt, priceEstimate.TotalPrice, priceEstimate.EstimatedDeliveryTimeInMinutes)
	if err != nil {
		return err
	}
	return nil
}
