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
	UserService     services.UserService
	PurchaseService services.PurchaseService
}

func NewUserController(userService services.UserService,
	purchaseService services.PurchaseService) UserController {
	return UserController{
		UserService:     userService,
		PurchaseService: purchaseService,
	}
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

type createPhonebookForm struct {
	Name string `json:"name" form:"name"`
}

type createTemplateForm struct {
	Expression string `json:"expression"`
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
		return c.String(http.StatusInternalServerError, "User not found!")
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

	err = u.UserService.AddContact(user, contact)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to add contact")
	}
	return c.String(http.StatusCreated, "Contact added successfully")
}

func (u *UserController) DeleteContact(c echo.Context) error {
	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	contactID, err := strconv.Atoi(c.Param("contactID"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid contact ID!")
	}

	contactExists := false
	for _, contact := range user.Contacts {
		if contact.ID == uint(contactID) {
			contactExists = true
		}
	}
	if !contactExists {
		return c.String(http.StatusBadRequest, "User has no contact with this id!")
	}

	err = u.UserService.DeleteContact(uint(contactID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete contact: "+err.Error())
	}
	return c.String(http.StatusOK, "Contact deleted successfully")
}

func (u *UserController) CreatePhoneBook(c echo.Context) error {
	form := new(createPhonebookForm)
	if err := c.Bind(form); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body!")
	}

	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	phoneBook := models.PhoneBook{
		UserID: user.ID,
		Name:   form.Name,
	}

	err = u.UserService.CreatePhoneBook(user, phoneBook)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create phone book: "+err.Error())
	}
	return c.String(http.StatusCreated, "Phone book created successfully")
}

func (u *UserController) AddNumberToPhoneBook(c echo.Context) error {

	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	phonebookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid PhoneBook ID")
	}
	numberID, err := strconv.Atoi(c.Param("nid"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Number ID")
	}

	idExists := false
	for _, book := range user.PhoneBooks {
		if int(book.ID) == phonebookID {
			idExists = true
			break
		}
	}
	if !idExists {
		return c.String(http.StatusBadRequest, "User has no phonebook with this id!")
	}

	phonebook, err := u.UserService.GetPhoneBook(uint(phonebookID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve phonebook with this id: "+err.Error())
	}

	number, err := u.UserService.GetNumberByID(uint(numberID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to retrieve number with this id: "+err.Error())
	}

	if err = u.UserService.UpdatePhoneBook(phonebook, number); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update phonebook: "+err.Error())
	}
	return c.String(http.StatusOK, "Number added successfully!")
}

func (u *UserController) DeletePhoneBook(c echo.Context) error {
	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	phonebookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid phonebook ID")
	}

	phonebookExists := false
	for _, phoneBook := range user.PhoneBooks {
		if phoneBook.ID == uint(phonebookID) {
			phonebookExists = true
		}
	}
	if !phonebookExists {
		return c.String(http.StatusNotAcceptable, "User has no phonebook with this id!")
	}

	err = u.UserService.DeletePhoneBook(uint(phonebookID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete phonebook: "+err.Error())
	}
	return c.String(http.StatusOK, "Phonebook deleted successfully")
}

func (u *UserController) AddTemplate(c echo.Context) error {
	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	body := createTemplateForm{}
	err = c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	template := &models.Template{
		UserID:     user.ID,
		Expression: body.Expression,
	}

	err = u.UserService.CreateTemplate(template)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, "Template created successfully!")
}

func (u *UserController) DeleteTemplate(c echo.Context) error {
	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	templateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid template ID")
	}

	templateExists := false
	for _, template := range user.Templates {
		if template.ID == uint(templateID) {
			templateExists = true
		}
	}
	if !templateExists {
		return c.String(http.StatusNotAcceptable, "User has no template with this id!")
	}

	err = u.UserService.DeleteTemplate(uint(templateID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete template: "+err.Error())
	}
	return c.String(http.StatusOK, "Template deleted successfully")
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
	return c.String(http.StatusOK, "Successfully updated user main number.")
}

func (u *UserController) BuyNumber(c echo.Context) error {
	numberId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Provide a valid number id!")
	}

	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	err = u.PurchaseService.BuyNumber(user, uint(numberId))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Successfully purchased number.")
}

func (u *UserController) PlaceRent(c echo.Context) error {
	numberId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Provide a valid number id!")
	}

	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	err = u.PurchaseService.PlaceRent(user, uint(numberId))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to place rent: "+err.Error())
	}
	return c.String(http.StatusOK, "Successfully placed rent.")
}

func (u *UserController) DropRent(c echo.Context) error {
	rentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Provide a valid rent id!")
	}

	user, err := u.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	err = u.PurchaseService.DropRent(user, uint(rentId))
	if err != nil {
		return c.String(http.StatusNotAcceptable, "Failed to drop rent: "+err.Error())
	}
	return c.String(http.StatusOK, "Successfully dropped rent.")
}

func validatePassword(password string) bool {
	//Constraints
	lengthConstraint := len(password) >= 8

	return lengthConstraint
}
