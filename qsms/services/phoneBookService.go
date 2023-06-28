package services

import (
	"qsms/models"
	"qsms/repository"
)

type PhoneBookService interface {
	CreatePhoneBook(phonebook *models.PhoneBook) error
	GetPhoneBook(phonebookId uint) (*models.PhoneBook, error)
	UpdatePhoneBook(phonebook *models.PhoneBook) error
	DeletePhoneBook(phonebookId uint) error
	AddContact(phonebook *models.PhoneBook, contact models.Contact) error
	DeleteContact(phonebook *models.PhoneBook, contactId uint) error
	GetContact(contactId uint) (*models.Contact, error)
	UpdateContact(phonebook *models.PhoneBook, contact *models.Contact) error
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

func (pbs *phoneBookService) AddContact(phonebook *models.PhoneBook, contact models.Contact) error {
	return pbs.phoneBookRepository.AddContact(phonebook, contact)
}

func (pbs *phoneBookService) DeleteContact(phonebook *models.PhoneBook, contactId uint) error {
	return pbs.phoneBookRepository.DeleteContact(phonebook, contactId)
}

func (pbs *phoneBookService) GetContact(contactId uint) (*models.Contact, error) {
	return pbs.phoneBookRepository.GetContact(contactId)
}

func (pbs *phoneBookService) UpdateContact(phonebook *models.PhoneBook, contact *models.Contact) error {
	return pbs.phoneBookRepository.UpdateContact(phonebook, contact)
}
