package controllers

import (
	"net/http"
	"qsms/models"
	"qsms/services"
	"qsms/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PhoneBookController struct {
	PhoneBookService services.PhoneBookService
	UserService      services.UserService
}

type createPhoneBookForm struct {
	Name string `json:"name" form:"name"`
}

func (pb *PhoneBookController) CreatePhoneBook(c echo.Context) error {
	form := new(createPhoneBookForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	user, err := pb.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	phoneBook := models.PhoneBook{
		UserID: user.ID,
		Name:   form.Name,
	}

	err = pb.PhoneBookService.CreatePhoneBook(user, &phoneBook)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create phone book")
	}

	return c.JSON(http.StatusCreated, "Phone book created successfully")
}

func (pb *PhoneBookController) GetPhoneBookByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	result, err := pb.PhoneBookService.GetPhoneBook(uint(id))
	if err != nil {
		return c.String(http.StatusNotFound, "PhoneBook Not Found!!")

	}
	return c.JSON(http.StatusOK, result)
}

func (pb *PhoneBookController) AddNumberToPhoneBook(c echo.Context) error {

	user, err := pb.UserService.UserByToken(utils.GetToken(c))
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

	phonebook, err := pb.PhoneBookService.GetPhoneBook(uint(phonebookID))
	if err != nil {
		return err
	}

	number, err := pb.PhoneBookService.GetNumberByID(uint(numberID))
	if err != nil {
		return err
	}

	if err = pb.PhoneBookService.UpdatePhoneBook(phonebook, number); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update phonebook"})
	}

	return c.String(http.StatusOK, "Number added successfully!")
}

func (pb *PhoneBookController) DeletePhoneBook(c echo.Context) error {
	phonebookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	_, err = pb.PhoneBookService.GetPhoneBook(uint(phonebookID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Phonebook not found"})
	}

	err = pb.PhoneBookService.DeletePhoneBook(uint(phonebookID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete phonebook"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Phonebook deleted successfully"})
}
