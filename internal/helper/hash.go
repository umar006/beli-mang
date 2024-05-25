package helper

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

var bcryptSalt = os.Getenv("BCRYPT_SALT")

func HashPassword(password string) ([]byte, error) {
	fmt.Println("HALO", bcryptSalt)
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
