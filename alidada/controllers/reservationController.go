package controllers

import (
	"alidada/models"
	"alidada/services"
	"alidada/utils"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"io/ioutil"
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

	order, err := rc.ReservationService.Reserve(reservationReq.PassengerIDs, reservationReq.FlightClassID)

	if err != nil {
		return c.String(http.StatusNotAcceptable, err.Error())
	}

	paymentRequest := models.Payment{
		MerchantID:  utils.ENV("MERCHANT_ID"),
		Amount:      int(order.Price),
		Description: "Reservation payment!",
		CallbackURL: "http://localhost:3000/api/reservation/verify",
	}

	jsonData, err := json.Marshal(&paymentRequest)
	if err != nil {
		return echo.ErrInternalServerError
	}

	request, err := http.NewRequest("POST", "https://sandbox.zarinpal.com/pg/v4/payment/request.json", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return echo.ErrInternalServerError
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return echo.ErrInternalServerError
	}

	var paymentResponse models.PaymentResponse
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if paymentResponse.Data.Code != 100 {
		return c.String(http.StatusBadGateway, "Failed to send payment!")
	}

	err = rc.ReservationService.SetAuthorityPair(paymentResponse.Data.Authority, order.ID)
	if err != nil {
		return err
	}

	err = c.Redirect(http.StatusOK, "https://sandbox.zarinpal.com/pg/StartPay/"+paymentResponse.Data.Authority)
	if err != nil {
		return c.String(http.StatusBadGateway, "Failed to redirect to payment page!")
	}

	return err
}

func (rc *ReservationController) Verify(c echo.Context) error {

	authority := c.Param("Authority")
	status := c.Param("Status")

	switch status {
	case "OK":

		order, err := rc.ReservationService.GetOrderByAuthority(authority)
		if err != nil {
			return errors.New("failed to find order by authority")
		}

		verifyRequest := models.VerifyRequest{
			MerchantID: utils.ENV("MERCHANT_ID"),
			Amount:     int(order.Price),
			Authority:  authority,
		}

		jsonData, err := json.Marshal(&verifyRequest)
		if err != nil {
			return echo.ErrInternalServerError
		}

		request, err := http.NewRequest("POST", "https://sandbox.zarinpal.com/pg/v4/payment/verify.json", bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return echo.ErrInternalServerError
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return echo.ErrInternalServerError
		}

		var verifyResponse models.VerifyResponse
		err = json.Unmarshal(body, &verifyResponse)

		if verifyResponse.Data.Code != 100 {
			return c.String(http.StatusNotAcceptable, "Payment not confirmed!")
		}

		err = rc.ReservationService.ConfirmOrder(order.ID, verifyResponse.Data.RefID)

	case "NOK":
		return c.String(http.StatusNotAcceptable, "Payment not confirmed!")
	default:
		return echo.ErrInternalServerError
	}

	return nil
}
