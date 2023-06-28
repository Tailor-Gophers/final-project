package controllers

import (
	"net/http"
	"qsms/services"

	"github.com/labstack/echo/v4"
)

type SMSController struct {
	PhoneBookService services.PhoneBookService
}

func (sms *SMSController) SendSMSToPhoneBooks(c echo.Context) error {
	return c.JSON(http.StatusOK, "SMS sent successfully")
}

func (sms *SMSController) SendSMSToPhoneNumbers(c echo.Context) error {
	return c.JSON(http.StatusOK, "SMS sent successfully")
}
