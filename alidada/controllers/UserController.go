package controllers

import (
	"final-project/alidada/models"
	"final-project/alidada/services"
	"final-project/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController() UserController {
	return UserController{UserService: services.NewUserService()}
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
	if !validatePassword(signupReq.Password) {
		return c.String(http.StatusBadRequest, "Password: "+signupReq.Password+" is too short.")
	}

	newUser := models.User{
		Email:       signupReq.Email,
		UserName:    signupReq.UserName,
		Password:    signupReq.Password,
		FirstName:   signupReq.FirstName,
		LastName:    signupReq.LastName,
		PhoneNumber: signupReq.PhoneNumber,
	}

	err = u.UserService.CreateUser(&newUser)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.String(http.StatusCreated, "Registration was successful")
}

func validatePassword(password string) bool {
	//Constraints
	lengthConstraint := len(password) >= 8
	//todo more constraints

	return lengthConstraint
}

type loginReq struct {
	// Email or User name required
	Email    string `json:"email"`     //opitional
	UserName string `json:"user_name"` //opitional
	Password string `json:"password"`
}

func (u *UserController) Login(c echo.Context) error {
	loginReq := &loginReq{}
	c.Bind(&loginReq)
	var user *models.User
	var err error
	if loginReq.UserName == "" {
		user, err = u.UserService.GetUserByEmail(loginReq.Email)
		fmt.Println(user)
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

func (u *UserController) GetUserByToken(c echo.Context) error {

	user, err := u.UserByToken(c)
	if err != nil {
		return echo.ErrUnauthorized
	} else {
		return c.JSON(http.StatusOK, user)

	}

}

func (u *UserController) LogOut(c echo.Context) error {

	token := GetToken(c)
	err := u.UserService.LogOut(token)
	if err != nil {
		return echo.ErrUnauthorized
	} else {
		return c.String(http.StatusOK, "Logout was successful")

	}

}

func GetToken(c echo.Context) string {

	authorization := c.Request().Header["Authorization"]
	Bearer := authorization[0]
	token := strings.Split(Bearer, "Bearer ")[1]
	return token
}

func (u *UserController) UserByToken(c echo.Context) (*models.User, error) {
	token := GetToken(c)
	user, err := u.UserService.UserByToken(token)
	if err != nil {
		return nil, echo.ErrInternalServerError
	}
	return user, nil
}
