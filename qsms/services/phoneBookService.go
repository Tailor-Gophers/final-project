package services

import (
	"fmt"
	"qsms/models"
	"qsms/repository"
)

type PhoneBookService interface {
	CreatePhoneBook(phonebook *models.PhoneBook) error
	GetPhoneBook(phonebookId uint) (*models.PhoneBook, error)
	UpdatePhoneBook(phonebook *models.PhoneBook) error
	DeletePhoneBook(phonebookId uint) error
	SendSMS(phoneNumber string, message string) error
}

type phoneBookService struct {
	phoneBookRepository repository.PhoneBookRepository
}

func NewPhoneBookService(repository repository.PhoneBookRepository) PhoneBookService {
	return &phoneBookService{
		phoneBookRepository: repository,
	}
}

func (pbs *phoneBookService) CreatePhoneBook(phonebook *models.PhoneBook) error {
	return pbs.phoneBookRepository.CreatePhoneBook(phonebook)
}

func (pbs *phoneBookService) GetPhoneBook(phonebookId uint) (*models.PhoneBook, error) {
	return pbs.phoneBookRepository.GetPhoneBook(phonebookId)
}

func (pbs *phoneBookService) UpdatePhoneBook(phonebook *models.PhoneBook) error {
	return pbs.phoneBookRepository.UpdatePhoneBook(phonebook)
}

func (pbs *phoneBookService) DeletePhoneBook(phonebookId uint) error {
	return pbs.phoneBookRepository.DeletePhoneBook(phonebookId)
}

func (pbs *phoneBookService) SendSMS(phoneNumber string, message string) error {
	// Implementation to send an SMS to the specified phone number with the given message

	// Example implementation:
	fmt.Printf("Sending SMS to %s: %s\n", phoneNumber, message)
	return nil
}
