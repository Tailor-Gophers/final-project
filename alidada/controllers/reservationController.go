package controllers

import (
	"alidada/services"
	"alidada/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReservationController struct {
	ReservationService services.ReservationService
	UserService        services.UserService
}

func NewReservationController() ReservationController {
	return ReservationController{
		ReservationService: services.NewReservationService(),
		UserService:        services.NewUserService(),
	}
}

type reservationForm struct {
	FlightClassID uint
	PassengerIDs  []uint
}

func (rc *ReservationController) Reserve(c echo.Context) error {
	reservationReq := &reservationForm{}
	err := c.Bind(reservationReq)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid reservation form!")
	}

	user, err := rc.UserService.UserByToken(utils.GetToken(c))

	//Check for passengers
	for _, id := range reservationReq.PassengerIDs {
		exists := false
		for _, passenger := range user.Passengers {
			if id == passenger.ID {
				exists = true
			}
			if !exists {
				return c.String(http.StatusBadRequest, "Passenger not found in user's passengers")
			}
		}
	}

	err = rc.ReservationService.Reserve(reservationReq.PassengerIDs, reservationReq.FlightClassID)

	if err != nil {
		return c.String(http.StatusNotAcceptable, err.Error())
	}
	return err
}
