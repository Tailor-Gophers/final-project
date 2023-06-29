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

func (pb *PhoneBookController) UpdatePhoneBook(c echo.Context) error {
	form := new(createPhoneBookForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	phonebookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	phonebook, err := pb.PhoneBookService.GetPhoneBook(uint(phonebookID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Phonebook not found"})
	}

	phonebook.Name = form.Name

	if err := pb.PhoneBookService.UpdatePhoneBook(phonebook); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update phonebook"})
	}

	return c.JSON(http.StatusOK, phonebook)
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
