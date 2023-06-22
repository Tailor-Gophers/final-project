package app

import (
	"mockapi/controllers"
	"mockapi/repository"
	"mockapi/services"

	"github.com/labstack/echo/v4"
	_ "gorm.io/driver/mysql"
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
	e.POST("/flights/:id/:class/reserve", FlightController.ReserveFlightCapacity)
	e.POST("/flights/:id/:class/return", FlightController.ReturnFlightCapacity)
	e.GET("/flights/filter/:airline/:aircraft/:departure", FlightController.GetFlightByFilter)
	e.GET("/flights/sort/:order", FlightController.GetFlightBySort)
	e.GET("/flight_class/:id", FlightController.GetFlightClassByID)
}
