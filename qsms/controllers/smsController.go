package controllers

import (
	"net/http"
	"qsms/services"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type SMSController struct {
	PhoneBookService services.PhoneBookService
	UserService      services.UserService
}

func (sms *SMSController) SendSMSToPhoneBooks(c echo.Context) error {
	phoneBookIDs := c.Param("phoneBookIDs")

	for _, phoneBookID := range strings.Split(phoneBookIDs, ",") {
		id, err := strconv.Atoi(phoneBookID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid phone book ID")
		}

		phoneBook, err := sms.PhoneBookService.GetPhoneBook(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error while retrieving phone book")
		}

		for _, numbers := range phoneBook.Numbers {
			phoneNumber := numbers.PhoneNumber

			err := sms.PhoneBookService.SendSMS(phoneNumber, "Your SMS message content")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, "Error while sending SMS")
			}
		}
	}
	return c.JSON(http.StatusOK, "SMS sent successfully")
}

func (sms *SMSController) SendSMSToPhoneNumbers(c echo.Context) error {
	phoneNumber := c.FormValue("number")
	message := c.FormValue("message")
	_, _ = phoneNumber, message

	// TODO: Implement the logic to send an SMS to the specified phone number using phoneNumber and message

	// Return a JSON response indicating success
	return c.JSON(http.StatusOK, "SMS sent successfully")
}
