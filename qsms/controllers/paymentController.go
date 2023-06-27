package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"qsms/models"
	"qsms/services"
	"qsms/utils"
	"strconv"
)

type PaymentController struct {
	UserService    services.UserService
	PaymentService services.PaymentService
}

func NewPaymentController(userService services.UserService, paymentService services.PaymentService) PaymentController {
	return PaymentController{
		UserService:    userService,
		PaymentService: paymentService,
	}
}

const (
	CallBackURL     = "http://localhost:3000/sms/payment/verify"
	PaymentRedirect = "https://sandbox.zarinpal.com/pg/StartPay/"
)

func (pc *PaymentController) AddBalance(c echo.Context) error {

	amount, err := strconv.Atoi(c.Param("amount"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid amount!")
	}

	user, err := pc.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return echo.ErrUnauthorized
	}

	code, authority, err := utils.RequestNewPayment(amount, CallBackURL)

	fmt.Println(code)
	fmt.Println(err.Error())

	if code != 100 {
		return c.String(http.StatusInternalServerError, "Failed to request payment gateway.")
	}

	transaction := &models.Transaction{
		UserID:    user.ID,
		Authority: authority,
		Amount:    amount,
		Confirmed: false,
	}

	err = pc.PaymentService.CreateTransaction(transaction)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create transaction!")
	}

	err = c.Redirect(http.StatusOK, PaymentRedirect+authority)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to redirect to payment!")
	}
	return c.String(http.StatusOK, "Redirected to payment page.")
}

func (pc *PaymentController) VerifyPayment(c echo.Context) error {

	authority := c.QueryParam("Authority")
	status := c.QueryParam("Status")

	if status == "NOK" {
		return c.String(http.StatusNotAcceptable, "Transaction not successfully!")
	}

	transaction, err := pc.PaymentService.GetTransaction(authority)
	if err != nil {
		return c.String(http.StatusNotAcceptable, "Failed to retrieve transaction!")
	}

	code, refId, err := utils.VerifyPayment(transaction.Amount, transaction.Authority)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if code != 100 {
		return c.String(http.StatusNotAcceptable, "Payment not confirmed!")
	}

	err = pc.PaymentService.ConfirmTransaction(&transaction, refId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to confirm payment!")
	}

	err = pc.UserService.AddBalance(transaction.UserID, transaction.Amount)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to add balance!")
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"ref_id":  strconv.Itoa(refId),
		"amount":  strconv.Itoa(transaction.Amount),
		"user_id": strconv.Itoa(int(transaction.UserID)),
	})
}
