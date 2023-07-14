package controllers

import (
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

func (pc *PaymentController) AddBalance(c echo.Context) error {

	amount, err := strconv.Atoi(c.Param("amount"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid amount!")
	}

	user, err := pc.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return echo.ErrUnauthorized
	}

	sandbox, _ := strconv.ParseBool(utils.ENV("SANDBOX"))
	zarinPay, err := utils.NewZarinpal(utils.ENV("MERCHANT_ID"), sandbox)
	if err != nil {
		return c.String(500, err.Error())
	}

	paymentURL, authority, statusCode, err := zarinPay.NewPaymentRequest(amount, utils.ENV("APP_URL")+"/sms/payment/verify", "User balance charge.", user.Email, "")

	if err != nil {
		if statusCode == -3 {
			return c.String(500, "Amount is not accepted in banking system")
		}
	}

	if statusCode != 100 {
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

	//err = c.Redirect(http.StatusFound, paymentURL)
	//if err != nil {
	//	return c.String(http.StatusInternalServerError, "Failed to redirect to payment!")
	//}

	return c.String(http.StatusOK, paymentURL)
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

	sandbox, _ := strconv.ParseBool(utils.ENV("SANDBOX"))
	zarinPay, err := utils.NewZarinpal(utils.ENV("MERCHANT_ID"), sandbox)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to initialize payment gateway: "+err.Error())
	}

	verified, refID, statusCode, err := zarinPay.PaymentVerification(transaction.Amount, authority)
	if err != nil {
		if statusCode == 101 {
			return c.String(500, "Payment is already verified")
		}
		return c.String(http.StatusNotAcceptable, "Payment was failed ")
	}

	refIDint, err := strconv.Atoi(refID)
	if err != nil {
		return err
	}

	if verified {
		err = pc.PaymentService.ConfirmTransaction(&transaction, refIDint)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to confirm payment!")
		}

		err = pc.UserService.AddBalance(transaction.UserID, transaction.Amount)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to add balance!")
		}
	} else {
		return c.String(http.StatusOK, "Payment cancelled ")
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"ref_id":  refID,
		"amount":  strconv.Itoa(transaction.Amount),
		"user_id": strconv.Itoa(int(transaction.UserID)),
	})
}
