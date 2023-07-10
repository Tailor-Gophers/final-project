package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"qsms/models"
	"qsms/services"
	"qsms/utils"
	"regexp"
)

type MessageController struct {
	UserService    services.UserService
	MessageService services.MessageService
}

func NewMessageController(UserService services.UserService, MessageService services.MessageService) MessageController {
	return MessageController{
		UserService:    UserService,
		MessageService: MessageService,
	}
}

type simpleRequestForm struct {
	Numbers      []string `json:"numbers"`
	ContactIDs   []int    `json:"contact_ids"`
	PhonebookIDs []int    `json:"phonebook_ids"`
	Text         string   `json:"text"`
	TemplateID   int      `json:"template_id"`
}

type periodicRequestForm struct {
	Numbers      []string `json:"numbers"`
	ContactIDs   []int    `json:"contact_ids"`
	PhonebookIDs []int    `json:"phonebook_ids"`
	Text         string   `json:"text"`
	TemplateID   int      `json:"template_id"`
	Interval     string   `json:"interval"`
}

func (sms *MessageController) SingleMessage(c echo.Context) error {
	user, err := sms.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	body := simpleRequestForm{}
	err = c.Bind(&body)

	if user.MainNumberID == 0 {
		return c.String(http.StatusNotAcceptable, "User has no main number set!")
	}

	if body.Text != "" {
		for _, number := range body.Numbers {
			err = sms.MessageService.SendSimpleMessage(user, number, body.Text)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, contactId := range body.ContactIDs {
			contact, err := getUserContact(user, uint(contactId))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			err = sms.MessageService.SendSimpleMessage(user, contact.PhoneNumber, body.Text)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, phonebookID := range body.PhonebookIDs {
			phoneBook, err := getUserPhonebook(user, uint(phonebookID))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			for _, number := range phoneBook.Numbers {
				err = sms.MessageService.SendSimpleMessage(user, number.PhoneNumber, body.Text)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
				}
			}
		}
	} else if body.TemplateID != 0 {

		template, err := sms.UserService.GetTemplate(uint(body.TemplateID))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get template: "+err.Error())
		}

		if template.UserID != user.ID {
			return c.String(http.StatusNotAcceptable, "User has no template with this id!")
		}

		for _, number := range body.Numbers {
			err = sms.MessageService.SendTemplateMessage(user, number, template.Expression)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, contactId := range body.ContactIDs {
			contact, err := getUserContact(user, uint(contactId))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			err = sms.MessageService.SendTemplateMessage(user, contact.PhoneNumber, template.Expression)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, phonebookID := range body.PhonebookIDs {
			phoneBook, err := getUserPhonebook(user, uint(phonebookID))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			for _, number := range phoneBook.Numbers {
				err = sms.MessageService.SendTemplateMessage(user, number.PhoneNumber, template.Expression)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
				}
			}
		}
	}
	return nil
}

func (sms *MessageController) PeriodicMessage(c echo.Context) error {
	user, err := sms.UserService.UserByToken(utils.GetToken(c))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized!")
	}

	body := periodicRequestForm{}
	err = c.Bind(&body)

	if user.MainNumberID == 0 {
		return c.String(http.StatusNotAcceptable, "User has no main number set!")
	}

	pattern := `^(?:[1-9]|[1-9][0-9]|100)[mhdM]$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(body.Interval) {
		return c.String(http.StatusBadRequest, "Invalid interval form.")
	}

	if body.Text != "" {
		for _, number := range body.Numbers {
			err = sms.MessageService.SendPeriodicSimpleMessage(user, number, body.Text, body.Interval)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, contactId := range body.ContactIDs {
			contact, err := getUserContact(user, uint(contactId))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			err = sms.MessageService.SendPeriodicSimpleMessage(user, contact.PhoneNumber, body.Text, body.Interval)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, phonebookID := range body.PhonebookIDs {
			phoneBook, err := getUserPhonebook(user, uint(phonebookID))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			for _, number := range phoneBook.Numbers {
				err = sms.MessageService.SendPeriodicSimpleMessage(user, number.PhoneNumber, body.Text, body.Interval)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
				}
			}
		}
	} else if body.TemplateID != 0 {

		template, err := sms.UserService.GetTemplate(uint(body.TemplateID))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get template: "+err.Error())
		}

		if template.UserID != user.ID {
			return c.String(http.StatusNotAcceptable, "User has no template with this id!")
		}

		for _, number := range body.Numbers {
			err = sms.MessageService.SendPeriodicTemplateMessage(user, number, template.Expression, body.Interval)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, contactId := range body.ContactIDs {
			contact, err := getUserContact(user, uint(contactId))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			err = sms.MessageService.SendPeriodicSimpleMessage(user, contact.PhoneNumber, template.Expression, body.Interval)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
			}
		}
		for _, phonebookID := range body.PhonebookIDs {
			phoneBook, err := getUserPhonebook(user, uint(phonebookID))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			for _, number := range phoneBook.Numbers {
				err = sms.MessageService.SendPeriodicTemplateMessage(user, number.PhoneNumber, template.Expression, body.Interval)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to send message: "+err.Error())
				}
			}
		}
	}
	return nil
}

func getUserContact(user *models.User, contactId uint) (models.Contact, error) {
	for _, contact := range user.Contacts {
		if contact.ID == contactId {
			return contact, nil
		}
	}
	return models.Contact{}, errors.New("contact not found for user")
}
func getUserPhonebook(user *models.User, phonebookId uint) (models.PhoneBook, error) {
	for _, phoneBook := range user.PhoneBooks {
		if phoneBook.ID == phonebookId {
			return phoneBook, nil
		}
	}
	return models.PhoneBook{}, errors.New("phonebook not found for user")
}
