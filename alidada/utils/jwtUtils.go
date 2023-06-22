package utils

import (
	"alidada/models"
	"github.com/labstack/echo/v4"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateTokenPair(user *models.User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["name"] = user.UserName
	//todo
	//claims["admin"] = user.Admin
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
