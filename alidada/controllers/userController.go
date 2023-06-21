package controllers

import (
	"alidada/models"
	"alidada/services"
	"alidada/utils"
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

type loginForm struct {
	// Email or Username required
	Email    string `json:"email"`     //optional
	UserName string `json:"user_name"` //optional
	Password string `json:"password"`
}

type passengerForm struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Gender         string `json:"gender"`
	DateOfBirth    string `json:"date_of_birth"`
	Nationality    string `json:"nationality"`
	PassportNumber string `json:"passport_number"`
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
		Passengers:  []models.Passenger{},
		PhoneNumber: signupReq.PhoneNumber,
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

func (u *UserController) UserByToken(c echo.Context) (*models.User, error) {
	token := GetToken(c)
	user, err := u.UserService.UserByToken(token)
	if err != nil {
		return nil, echo.ErrInternalServerError
	}
	return user, nil
}

func (u *UserController) CreatePassenger(c echo.Context) error {

	passengerReq := &passengerForm{}
	err := c.Bind(passengerReq)
	if err != nil {
		return c.String(http.StatusBadRequest, "All passenger fields must be provided!")
	}
	newPassenger := &models.Passenger{
		FirstName:      passengerReq.FirstName,
		LastName:       passengerReq.LastName,
		Gender:         passengerReq.Gender,
		DateOfBirth:    passengerReq.DateOfBirth,
		Nationality:    passengerReq.Nationality,
		PassportNumber: passengerReq.PassportNumber,
	}

	user, err := u.UserByToken(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "You must be logged in!")
	}
	return u.UserService.CreatePassenger(user, newPassenger)
}

func (u *UserController) GetPassengers(c echo.Context) error {
	user, err := u.UserByToken(c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "You must be logged in!")
	}
	passengers, _ := u.UserService.GetPassengers(user)
	return c.JSON(http.StatusOK, passengers)
}

func GetToken(c echo.Context) string {

	authorization := c.Request().Header["Authorization"]
	Bearer := authorization[0]
	token := strings.Split(Bearer, "Bearer ")[1]
	return token
}
func validatePassword(password string) bool {
	//Constraints
	lengthConstraint := len(password) >= 8
	//todo more constraints

	return lengthConstraint
}
