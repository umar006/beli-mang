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
	CreateAdmin(ctx context.Context, body domain.RegisterRequest) (string, *fiber.Error)
	CreateCustomer(ctx context.Context, body domain.RegisterRequest) (string, *fiber.Error)
	Login(ctx context.Context, body domain.LoginRequest) (string, *fiber.Error)
	GetPriceEstimation(ctx context.Context, body domain.PriceEstimateRequest) (domain.PriceEstimateResponse, *fiber.Error)
}

type userService struct {
	db           *pgx.Conn
	userRepo     repository.UserRepo
	merchantRepo repository.MerchantRepo
	itemRepo     repository.MerchantItemRepo
}

func NewUser(db *pgx.Conn, userRepo repository.UserRepo, merchantRepo repository.MerchantRepo, itemRepo repository.MerchantItemRepo) UserService {
	return &userService{
		db:           db,
		userRepo:     userRepo,
		merchantRepo: merchantRepo,
		itemRepo:     itemRepo,
	}
}

func (us *userService) CreateAdmin(ctx context.Context, body domain.RegisterRequest) (string, *fiber.Error) {
	admin := body.NewUserFromDTO()
	admin.Role = domain.RoleAdmin

	hashedPassword, err := helper.HashPassword(admin.Password)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	admin.Password = string(hashedPassword)
	err = us.userRepo.Create(ctx, us.db, admin)
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

func (us *userService) CreateCustomer(ctx context.Context, body domain.RegisterRequest) (string, *fiber.Error) {
	customer := body.NewUserFromDTO()
	customer.Role = domain.RoleCustomer

	hashedPassword, err := helper.HashPassword(customer.Password)
	if err != nil {
		return "", domain.NewErrInternalServerError(err.Error())
	}

	customer.Password = string(hashedPassword)
	err = us.userRepo.Create(ctx, us.db, customer)
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

	token, err := helper.GenerateJWTToken(customer)
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

func (us *userService) GetPriceEstimation(ctx context.Context, body domain.PriceEstimateRequest) (domain.PriceEstimateResponse, *fiber.Error) {
	merchantIDs := []string{}

	for _, order := range body.Orders {
		merchantIDs = append(merchantIDs, order.MerchantID)
	}

	merchantList, err := us.merchantRepo.GetMerchantListByIDs(ctx, us.db, merchantIDs)
	if err != nil {
		return domain.PriceEstimateResponse{}, domain.NewErrInternalServerError(err.Error())
	}
	if len(merchantList) != len(merchantIDs) {
		return domain.PriceEstimateResponse{}, domain.NewErrNotFound("some merchants are not found")
	}

	return domain.PriceEstimateResponse{}, nil
}
