package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	PaymentRequestURL = "https://sandbox.zarinpal.com/pg/v4/payment/request.json"
	PaymentVerifyURL  = "https://sandbox.zarinpal.com/pg/v4/payment/verify.json"
)

type Payment struct {
	MerchantID  string   `json:"merchant_id"`
	Amount      int      `json:"amount"`
	CallbackURL string   `json:"callback_url"`
	Description string   `json:"description"`
	MetaData    MetaData `json:"metadata"`
}

type MetaData struct {
	Email  string
	Mobile string
}

type PaymentResponse struct {
	Data   Data     `json:"data"`
	Errors []string `json:"errors"`
}
type Data struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Authority string `json:"authority"`
	FeeType   string `json:"fee_type"`
	Fee       int    `json:"fee"`
}

type AuthorityPair struct {
	OrderID   uint
	Authority string `gorm:"size:255"`
}

type VerifyRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int    `json:"amount"`
	Authority  string `json:"authority"`
}
type VerifyResponse struct {
	Data   VerifyData `json:"data"`
	Errors []string   `json:"errors"`
}
type VerifyData struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CardHash string `json:"card_hash"`
	CardPan  string `json:"card_pan"`
	RefID    int    `json:"ref_id"`
	FeeType  string `json:"fee_type"`
	Fee      int    `json:"fee"`
}

func RequestNewPayment(amount int, callBack string) (int, string, error) {
	paymentRequest := Payment{
		MerchantID:  ENV("MERCHANT_ID"),
		Amount:      amount,
		Description: "Reservation payment!",
		CallbackURL: callBack,
	}

	jsonData, err := json.Marshal(&paymentRequest)
	if err != nil {
		return -1, "", err
	}

	fmt.Println(string(jsonData))

	request, err := http.NewRequest("POST", PaymentRequestURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return -2, "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -3, "", err
	}

	fmt.Println(string(body))

	var paymentResponse PaymentResponse
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return -4, "", err
	}

	if paymentResponse.Data.Code != 100 {
		return paymentResponse.Data.Code, "", err
	}

	return paymentResponse.Data.Code, paymentResponse.Data.Authority, nil
}

func VerifyPayment(amount int, authority string) (int, int, error) {
	verifyRequest := VerifyRequest{
		MerchantID: ENV("MERCHANT_ID"),
		Amount:     amount,
		Authority:  authority,
	}

	jsonData, err := json.Marshal(&verifyRequest)
	if err != nil {
		return -1, -1, err
	}

	request, err := http.NewRequest("POST", PaymentVerifyURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return -1, -1, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, -1, err
	}

	var verifyResponse VerifyResponse
	err = json.Unmarshal(body, &verifyResponse)

	if verifyResponse.Data.Code != 100 {
		return verifyResponse.Data.Code, -1, errors.New("transaction not completed")
	}

	return verifyResponse.Data.Code, verifyResponse.Data.RefID, nil
}
