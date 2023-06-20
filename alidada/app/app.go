package app

import (
	"alidada/controllers"

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
	flightController := controllers.NewFlightController()
	userController := controllers.NewUserController()
	authGroup := e.Group("/api/auth")
	authGroup.POST("/signup", userController.Signup)
	authGroup.POST("/login", userController.Login)
	authGroup.GET("/me", userController.GetUserByToken)
	authGroup.POST("/logout", userController.LogOut)

	//todo login
	userGroup := e.Group("/api/user")
	userGroup.GET("/passengers", userController.GetPassengers)
	userGroup.POST("/AddPassenger", userController.CreatePassenger)

	// authGroup.POST("/logout", userController.Login)

	// mockapi
	e.GET("/flights", flightController.SearchFlights)

}

func qsmsRouting(e *echo.Echo) {
}
