package helper

import (
	"beli-mang/internal/domain"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

var (
	bcryptSalt = os.Getenv("BCRYPT_SALT")
	jwtSecret  = os.Getenv("JWT_SECRET")
)

func HashPassword(password string) ([]byte, error) {
	parsedSalt, err := strconv.Atoi(bcryptSalt)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), parsedSalt)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func GenerateJWTToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
