package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/helper"
	"beli-mang/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserService interface {
	CreateAdmin(ctx context.Context, body domain.AdminRequest) (string, *fiber.Error)
	Login(ctx context.Context, body domain.LoginRequest) (string, *fiber.Error)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				column := helper.ExtractColumnFromPgErr(pgErr)
				return "", domain.NewErrConflict(fmt.Sprintf("%s already exists", column))
			}
		}
		return "", domain.NewErrInternalServerError(err.Error())
	}

	token, err := helper.GenerateJWTToken(admin)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	return token, nil
}

func (us *userService) Login(ctx context.Context, body domain.LoginRequest) (string, *fiber.Error) {
	user, err := us.userRepo.GetUserByUsername(ctx, us.db, body.Username)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	ok := helper.ComparePassword(user.Password, body.Password)
	if !ok {
		return "", domain.NewErrBadRequest("password is wrong")
	}

	token, err := helper.GenerateJWTToken(user)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	return token, nil
}
