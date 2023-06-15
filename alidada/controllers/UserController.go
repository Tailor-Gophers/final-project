package controllers

import (
	"final-project/alidada/models"
	"final-project/alidada/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserController struct {
	UserService services.UserService
}

type signUpForm struct {
	Email       string `json:"email"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

func (u *UserController) Signup(c echo.Context) error {
	signupReq := &signUpForm{}
	err := c.Bind(&signupReq)
	if err != nil {
		return c.String(http.StatusBadRequest, "All user fields must be provided!")
	}

	//Check username not taken
	_, err = u.UserService.GetUserByUserName(signupReq.UserName)
	if err == nil {
		return c.String(http.StatusBadRequest, "Username"+signupReq.UserName+" already exists.")
	}

	//Check email not taken
	_, err = u.UserService.GetUserByEmail(signupReq.Email)
	if err == nil {
		return c.String(http.StatusBadRequest, "Email "+signupReq.Email+" already exists.")
	}

	//Check password
	if validatePassword(signupReq.Password) {
		return c.String(http.StatusBadRequest, "Password: "+signupReq.Password+" is too short.")
	}

	newUser := models.User{
		Email:       signupReq.Email,
		UserName:    signupReq.UserName,
		Password:    signupReq.Password,
		FirstName:   signupReq.FirstName,
		LastName:    signupReq.LastName,
		PhoneNumber: signupReq.PhoneNumber,
		Admin:       false,
	}

	err = u.UserService.CreateUser(&newUser)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusCreated, nil)
}

func validatePassword(password string) bool {
	//Constraints
	lengthConstraint := len(password) >= 8
	//todo more constraints

	return lengthConstraint
}
