package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"qsms/models"
	"qsms/services"
	"qsms/utils"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(service services.UserService) UserController {
	return UserController{UserService: service}
}

type signUpForm struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type loginForm struct {
	Email    string `json:"email"`     //optional
	UserName string `json:"user_name"` //optional
	Password string `json:"password"`
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
	if !validatePassword(signupReq.Password) {
		return c.String(http.StatusBadRequest, "Password: "+signupReq.Password+" is too short.")
	}

	newUser := models.User{
		Email:    signupReq.Email,
		UserName: signupReq.UserName,
		Password: signupReq.Password,
		Balance:  0,
		Disable:  false,
		Admin:    false,
	}

	err = u.UserService.CreateUser(&newUser)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.String(http.StatusCreated, "Registration was successful")
}

func (u *UserController) Login(c echo.Context) error {
	loginReq := &loginForm{}
	err := c.Bind(&loginReq)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid login form.")
	}

	var user *models.User
	if loginReq.UserName == "" {
		user, err = u.UserService.GetUserByEmail(loginReq.Email)
	} else {
		user, err = u.UserService.GetUserByUserName(loginReq.UserName)
	}
	if err != nil {
		return echo.ErrInternalServerError
	}
	if utils.ValidatePassword(user.Password, loginReq.Password) {
		token, err := utils.GenerateTokenPair(user)
		if err != nil {
			return echo.ErrInternalServerError
		}
		err = u.UserService.SaveToken(user, token)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	}
	return echo.ErrUnauthorized
}

func (u *UserController) LogOut(c echo.Context) error {
	token := utils.GetToken(c)
	err := u.UserService.LogOut(token)
	if err != nil {
		return echo.ErrUnauthorized
	} else {
		return c.String(http.StatusOK, "Logout was successful")

	}

}

func (u *UserController) GetUserByToken(c echo.Context) error {

	user, err := u.UserByToken(c)
	if err != nil {
		return echo.ErrUnauthorized
	} else {
		return c.JSON(http.StatusOK, user)
	}
}

func (u *UserController) UserByToken(c echo.Context) (*models.User, error) {
	token := utils.GetToken(c)
	user, err := u.UserService.UserByToken(token)
	if err != nil {
		return nil, echo.ErrInternalServerError
	}
	return user, nil
}

func validatePassword(password string) bool {
	//Constraints
	lengthConstraint := len(password) >= 8

	return lengthConstraint
}
