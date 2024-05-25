package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type UserService interface {
	CreateAdmin(ctx context.Context, body domain.AdminRequest) *fiber.Error
}

type userService struct {
	db       *pgx.Conn
	userRepo repository.UserRepo
}

func NewUser(db *pgx.Conn, userRepo repository.UserRepo) UserService {
	return &userService{
		db:       db,
		userRepo: userRepo,
	}
}

func (us *userService) CreateAdmin(ctx context.Context, body domain.AdminRequest) *fiber.Error {
	admin := body.NewUserFromDTO()

	err := us.userRepo.CreateAdmin(ctx, us.db, admin)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
