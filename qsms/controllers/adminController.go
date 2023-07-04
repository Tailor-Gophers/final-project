package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"qsms/models"
	"qsms/services"
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

func (ac *AdminController) AddNumber(c echo.Context) error {
	body := numberRequestForm{}
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body!")
	}

	number := &models.Number{
		UserID:      0,
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
