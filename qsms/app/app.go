package app

import (
	"github.com/labstack/echo/v4"
	"qsms/controllers"
)

type App struct {
	E *echo.Echo
}

func NewApp() *App {
	e := echo.New()

	userController := controllers.NewUserController()

	userGroup := e.Group("/sms/user")
	userGroup.POST("/signup", userController.Signup)
	userGroup.POST("/login", userController.Login)
	userGroup.GET("/me", userController.GetUserByToken)
	userGroup.POST("/logout", userController.LogOut)

	return &App{
		E: e,
	}
}

func (a *App) Start(addr string) error {
	a.E.Logger.Fatal(a.E.Start(addr))
	return nil
}
