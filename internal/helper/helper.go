package helper

import (
	"beli-mang/internal/domain"
	"os"
	"regexp"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		"role":     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ExtractColumnFromPgErr(pgErr *pgconn.PgError) string {
	// Define the regular expression pattern
	// This pattern captures the key and the value separately
	re := regexp.MustCompile(`\((.*?)\)=\((.*?)\)`)

	// Find all matches
	matches := re.FindStringSubmatch(pgErr.Detail)

	// Check if both column and value were captured
	if len(matches) < 2 {
		return "no column found"
	}
	column := matches[1]
	value := matches[2]

	return column + " " + value
}
