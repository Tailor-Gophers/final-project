package middlewares

import (
	"alidada/controllers"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userController := controllers.NewUserController()
		user, err := userController.UserByToken(c)
		if err != nil {
			return c.String(http.StatusUnauthorized, "You must be logged in!")
		}
		if user.IsAdmin == false {
			return c.String(http.StatusUnauthorized, "You must be admin !")
		}

		return next(c)
	}
}
