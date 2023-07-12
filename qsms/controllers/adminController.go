package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"qsms/models"
	"qsms/services"
	"qsms/utils"
	"strconv"
)

type AdminController struct {
	AdminService services.AdminService
}

func NewAdminController(service services.AdminService) AdminController {
	return AdminController{AdminService: service}
}

type numberRequestForm struct {
	PhoneNumber string `json:"phone_number"`
	Price       int    `json:"price"`
}
type searchRequestForm struct {
	Words []string `json:"words"`
}
type setFeeRequestForm struct {
	SimpleMessageFee   int `json:"simple_message_fee"`
	PeriodicMessageFee int `json:"periodic_message_fee"`
	TemplateFee        int `json:"template_fee"`
}

func (ac *AdminController) AddNumber(c echo.Context) error {
	body := numberRequestForm{}
	err := c.Bind(&body)
	if err != nil || len(body.PhoneNumber) < 10 || body.Price < 0 {
		return c.String(http.StatusBadRequest, "Invalid request body!")
	}

	number := &models.Number{
		PhoneNumber: body.PhoneNumber,
		Price:       body.Price,
		Active:      false,
	}

	err = ac.AdminService.AddNumber(number)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to add number: "+err.Error())
	}
	return c.String(http.StatusOK, "Number added successfully!")
}

func (ac *AdminController) SuspendUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Provide a valid user id!")
	}

	err = ac.AdminService.SuspendUser(uint(userId))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to suspend user: "+err.Error())
	}
	return c.String(http.StatusOK, "User suspended successfully!")
}

func (ac *AdminController) UnSuspendUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Provide a valid user id!")
	}

	err = ac.AdminService.UnSuspendUser(uint(userId))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to unsuspend user: "+err.Error())
	}
	return c.String(http.StatusOK, "User unsuspended successfully!")
}

func (ac *AdminController) CountUserMessages(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Provide a valid user id!")
	}

	count, err := ac.AdminService.CountUserMessages(uint(userId))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]int{"user_id": userId, "message_count": count})
}

func (ac *AdminController) SearchMessages(c echo.Context) error {
	body := searchRequestForm{}
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body!")
	}

	messageFound, err := ac.AdminService.SearchMessages(body.Words)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to search messages: "+err.Error())
	}
	return c.JSON(http.StatusOK, messageFound)
}

func (ac *AdminController) AddBadWord(c echo.Context) error {
	body := searchRequestForm{}
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body!")
	}

	config := utils.LoadConfig()
	utils.AddBadWord(config, body.Words)

	err = utils.SaveConfig(config)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save config: "+err.Error())
	}
	return c.String(http.StatusOK, "Success!")
}

func (ac *AdminController) SetFee(c echo.Context) error {
	body := setFeeRequestForm{}
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid body!")
	}

	config := utils.LoadConfig()

	config.SimpleMessageFee = body.SimpleMessageFee
	config.PeriodicMessageFee = body.PeriodicMessageFee
	config.TemplateFee = body.TemplateFee

	err = utils.SaveConfig(config)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save config: "+err.Error())
	}
	return c.String(http.StatusOK, "Success!")
}
