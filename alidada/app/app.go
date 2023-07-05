package app

import (
	"alidada/controllers"

	echo "github.com/labstack/echo/v4"
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
	reservationController := controllers.NewReservationController()

	authGroup := e.Group("/api/auth")
	authGroup.POST("/signup", userController.Signup)
	authGroup.POST("/login", userController.Login)
	authGroup.GET("/me", userController.GetUserByToken)
	authGroup.POST("/logout", userController.LogOut)
	e.GET("/api/user/tickets", userController.GetMyTickets)
	e.POST("/api/user/tickets/cancell/:id", userController.CancellTicket)
	e.GET("/api/user/pdftickets/:id", userController.GetMyTicketsPdf)

	userGroup := e.Group("/api/user")
	userGroup.GET("/passengers", userController.GetPassengers)
	userGroup.POST("/AddPassenger", userController.CreatePassenger)

	reservationGroup := e.Group("/api/reservation")
	reservationGroup.POST("/reserve", reservationController.Reserve)
	reservationGroup.GET("/verify", reservationController.Verify) //http://www.yoursite.ir/?Authority=A00000000000000000000000000202690354&Status=OK ????

	e.GET("/flights/search", flightController.SearchFlightsDay)
	e.GET("/flights/sort", flightController.SearchFlightsSort)
	e.GET("/flights/filter", flightController.FiletrFlights)
}

func qsmsRouting(e *echo.Echo) {
}
