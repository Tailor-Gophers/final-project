package controllers

import (
	"net/http"
	"qsms/models"
	"qsms/services"
	"qsms/utils"
	"strconv"

	"github.com/labstack/echo/v4"
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

type createContactForm struct {
	Name        string `json:"name" form:"name"`
	PhoneNumber string `json:"phone"`
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

func (u *UserController) GetPhoneNumbersToBuy(c echo.Context) error {
	numbers, err := u.UserService.GetAvailablePhoneNumbers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve numbers: "+err.Error())
	}
	return c.JSON(http.StatusOK, numbers)
}

func (u *UserController) AddContact(c echo.Context) error {

	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	form := new(createContactForm)
	if err = c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	contact := &models.Contact{
		UserID:      user.ID,
		Name:        form.Name,
		PhoneNumber: form.PhoneNumber,
	}

	err = u.UserService.AddContact(contact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add contact"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Contact added successfully"})
}

func (u *UserController) DeleteContact(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	contactID, err := strconv.Atoi(c.Param("contactID"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	user, err := u.UserService.GetUserByID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	err = u.UserService.DeleteContact(user, uint(contactID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete contact"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact deleted successfully"})
}

func (u *UserController) UpdateContact(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	contactID, err := strconv.Atoi(c.Param("contactID"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	user, err := u.UserService.GetUserByID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	contact, err := u.UserService.GetContact(uint(contactID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contact not found"})
	}

	form := new(createContactForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	contact.Name = form.Name
	contact.PhoneNumber = form.PhoneNumber

	err = u.UserService.UpdateContact(user, contact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update contact"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact updated successfully"})
}

func (u *UserController) SetMainNumber(c echo.Context) error {

	numberId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Provide a valid number id!")
	}

	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	if int(user.MainNumberID) == numberId {
		return c.String(http.StatusNotAcceptable, "User's main number id already is "+strconv.Itoa(numberId))
	}

	numberIdExists := false
	for _, number := range user.Numbers {
		if int(number.ID) == numberId {
			numberIdExists = true
		}
	}
	if !numberIdExists {
		return c.String(http.StatusNotAcceptable, "User doesn't own a number with id "+strconv.Itoa(numberId))
	}

	err = u.UserService.SetMainNumber(user, uint(numberId))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func validatePassword(password string) bool {
	//Constraints
	lengthConstraint := len(password) >= 8

	return lengthConstraint
}
