package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"qsms/models"
	"strings"
	"time"
)

func HashPassword(pass string) (string, error) {
	if len(pass) == 0 {
		return "", errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash), err
}

func ValidatePassword(givenPass, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(givenPass), []byte(pass)) == nil
}

func HashToken(pass string) (string, error) {
	// if len(pass) == 0 {
	// 	return "", errors.New("password cannot be empty")
	// }

	// h := sha256.New()
	// h.Write([]byte(pass))
	// // Calculate and print the hash
	// s := fmt.Sprintf("%x", h.Sum(nil))
	return pass, nil
}

func ENV(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func GenerateTokenPair(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["name"] = user.UserName
	claims["admin"] = user.Admin
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GetToken(c echo.Context) string {
	authorization := c.Request().Header["Authorization"]
	Bearer := authorization[0]
	token := strings.Split(Bearer, "Bearer ")[1]
	return token
}
