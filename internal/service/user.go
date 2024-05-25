package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/helper"
	"beli-mang/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type UserService interface {
	CreateAdmin(ctx context.Context, body domain.AdminRequest) (string, *fiber.Error)
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

func (us *userService) CreateAdmin(ctx context.Context, body domain.AdminRequest) (string, *fiber.Error) {
	admin := body.NewUserFromDTO()

	hashedPassword, err := helper.HashPassword(admin.Password)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	admin.Password = string(hashedPassword)
	err = us.userRepo.CreateAdmin(ctx, us.db, admin)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	token, err := helper.GenerateJWTToken(admin)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	return token, nil
}
