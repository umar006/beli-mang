package service

import (
	"beli-mang/internal/domain"
	"beli-mang/internal/repository"
	"context"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var bcryptSalt = os.Getenv("BCRYPT_SALT")

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

	parsedSalt, err := strconv.Atoi(bcryptSalt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), parsedSalt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	admin.Password = string(hashedPassword)
	err = us.userRepo.CreateAdmin(ctx, us.db, admin)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
