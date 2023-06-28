package controllers

import (
	"net/http"
	"qsms/models"
	"qsms/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PhoneBookController struct {
	PhoneBookService services.PhoneBookService
}

type createPhoneBookForm struct {
	Name string `json:"name" form:"name"`
}

type createContactForm struct {
	Name        string `json:"name" form:"name"`
	PhoneNumber string `json:"phone"`
}

func (pb *PhoneBookController) CreatePhoneBook(c echo.Context) error {
	form := new(createPhoneBookForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	phoneBook := models.PhoneBook{
		Name: form.Name,
	}

	err := pb.PhoneBookService.CreatePhoneBook(&phoneBook)
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

func (pb *PhoneBookController) AddContact(c echo.Context) error {
	phonebookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	phonebook, err := pb.PhoneBookService.GetPhoneBook(uint(phonebookID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Phonebook not found"})
	}

	form := new(createContactForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	contact := models.Contact{
		Name:        form.Name,
		PhoneNumber: form.PhoneNumber,
	}

	err = pb.PhoneBookService.AddContact(phonebook, contact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add contact"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Contact added successfully"})
}

func (pb *PhoneBookController) DeleteContact(c echo.Context) error {
	phonebookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	contactID, err := strconv.Atoi(c.Param("contactID"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	phonebook, err := pb.PhoneBookService.GetPhoneBook(uint(phonebookID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Phonebook not found"})
	}

	err = pb.PhoneBookService.DeleteContact(phonebook, uint(contactID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete contact"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact deleted successfully"})
}

func (pb *PhoneBookController) UpdateContact(c echo.Context) error {
	phonebookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	contactID, err := strconv.Atoi(c.Param("contactID"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	phonebook, err := pb.PhoneBookService.GetPhoneBook(uint(phonebookID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Phonebook not found"})
	}

	contact, err := pb.PhoneBookService.GetContact(uint(contactID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contact not found"})
	}

	form := new(createContactForm)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	contact.Name = form.Name
	contact.PhoneNumber = form.PhoneNumber

	err = pb.PhoneBookService.UpdateContact(phonebook, contact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update contact"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact updated successfully"})
}
