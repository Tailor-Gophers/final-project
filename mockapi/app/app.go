package app

import (
	"github.com/labstack/echo"
	_ "gorm.io/driver/mysql"
	"mockapi/controllers"
	"mockapi/repository"
	"mockapi/services"
)

type App struct {
	E *echo.Echo
}

func NewApp() *App {
	e := echo.New()
	routing(e)
	return &App{
		E: e,
	}
}

func (a *App) Start(addr string) error {
	a.E.Logger.Fatal(a.E.Start(addr))
	return nil
}

func routing(e *echo.Echo) {
	flightRepo := repository.NewGormFlightRepository()
	FlightService := services.NewFlightService(flightRepo)
	FlightController := controllers.FlightController{FlightService: FlightService}

	// Public routes
	e.GET("/flights/:id", FlightController.GetFlightByID)
	e.GET("/flights/:origin/:destination/:date", FlightController.GetFlightByDate)
	e.GET("/flights/planes", FlightController.GetPlanesList)
	e.GET("/flights/cities", FlightController.GetCitiesList)
	e.GET("/flights/days", FlightController.GetDaysList)
	e.POST("/flights/:id/reserve", FlightController.ReserveFlightCapacity)
	e.POST("/flights/:id/return", FlightController.ReturnFlightCapacity)

}
