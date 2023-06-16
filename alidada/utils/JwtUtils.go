package utils

import (
	"alidada/models"
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
