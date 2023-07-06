package controllers

import (
	"alidada/services"
	"alidada/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
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
		}
		if !exists {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Passenger id %d not found in user's passengers", id))
		}
	}

	order, err := rc.ReservationService.Reserve(reservationReq.PassengerIDs, reservationReq.FlightClassID)

	if err != nil {
		return c.String(http.StatusNotAcceptable, err.Error())
	}
	sandbox, _ := strconv.ParseBool(utils.ENV("SANDBOX"))
	zarinPay, err := utils.NewZarinpal(utils.ENV("MERCHANT_ID"), sandbox)
	if err != nil {
		return c.String(500, err.Error())
	}

	paymentURL, authority, statusCode, err := zarinPay.NewPaymentRequest(int(order.Price), "http://localhost:3000/api/reservation/verify", fmt.Sprintf("order id %d", order.ID), user.Email, user.PhoneNumber)
	if err != nil {
		if statusCode == -3 {
			return c.String(500, "Amount is not accepted in banking system")
		}
		log.Fatal(err)
	}
	err = rc.ReservationService.SetAuthorityPair(authority, order.ID)
	if err != nil {
		return err
	}
	return c.String(200, paymentURL)
}

func (rc *ReservationController) Verify(c echo.Context) error {
	sandbox, _ := strconv.ParseBool(utils.ENV("SANDBOX"))

	zarinPay, err := utils.NewZarinpal(utils.ENV("MERCHANT_ID"), sandbox)
	if err != nil {
		log.Fatal(err)
	}
	authority := c.QueryParam("Authority")
	order, err := rc.ReservationService.GetOrderByAuthority(authority)
	if err != nil {
		return c.String(500, "failed to find order by authority")
	}
	amount := int(order.Price)

	verified, refID, statusCode, err := zarinPay.PaymentVerification(amount, authority)
	if err != nil {
		if statusCode == 101 {
			return c.String(500, "Payment is already verified")
		}
		return c.String(200, "payment was failed ")
	}

	refIDInt, _ := strconv.ParseInt(refID, 10, 0)

	if verified {
		err = rc.ReservationService.ConfirmOrder(order.ID, int(refIDInt))
		return c.String(200, "payment was successful ")

	} else {
		return c.String(200, "payment cancelled ")

	}

	return nil
}
