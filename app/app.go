package app

import (
	"final-project/alidada/controllers"
	"github.com/labstack/echo/v4"
)

type App struct {
	E *echo.Echo
}

func NewApp() *App {
	e := echo.New()
	alidadaRouting(e)
	qsmsRouting(e)
	return &App{
		E: e,
	}
}

func (a *App) Start(addr string) error {
	a.E.Logger.Fatal(a.E.Start(addr))
	return nil
}

func alidadaRouting(e *echo.Echo) {
	userController := controllers.NewUserController()
	authGroup := e.Group("/api/auth")
	authGroup.POST("/signup", userController.Signup)
}

func qsmsRouting(e *echo.Echo) {
}
