package repository

import (
	"qsms/db"
	"qsms/models"

	"gorm.io/gorm"
)

type PhoneBookRepository interface {
	CreatePhoneBook(phonebook *models.PhoneBook) error
	GetPhoneBook(phonebookId uint) (*models.PhoneBook, error)
	UpdatePhoneBook(phonebook *models.PhoneBook) error
	DeletePhoneBook(phonebookId uint) error
	AddContact(phonebook *models.PhoneBook, contact models.Contact) error
	DeleteContact(phonebook *models.PhoneBook, contactId uint) error
	GetContact(contactId uint) (*models.Contact, error)
	UpdateContact(phonebook *models.PhoneBook, contact *models.Contact) error
}

type phoneBookGormRepository struct {
	db *gorm.DB
}

func NewGormPhoneBookRepository() PhoneBookRepository {
	return &phoneBookGormRepository{
		db: db.GetDbConnection(),
	}
}

func (pb *phoneBookGormRepository) CreatePhoneBook(phonebook *models.PhoneBook) error {
	return pb.db.Create(phonebook).Error
}

func (pb *phoneBookGormRepository) GetPhoneBook(phonebookId uint) (*models.PhoneBook, error) {
	var phonebook models.PhoneBook
	err := pb.db.Preload("Contacts").First(&phonebook, phonebookId).Error
	if err != nil {
		return nil, err
	}
	return &phonebook, nil
}

func (pb *phoneBookGormRepository) UpdatePhoneBook(phonebook *models.PhoneBook) error {
	return pb.db.Save(phonebook).Error
}

func (pb *phoneBookGormRepository) DeletePhoneBook(phonebookId uint) error {
	return pb.db.Delete(&models.PhoneBook{}, phonebookId).Error
}

func (pb *phoneBookGormRepository) AddContact(phonebook *models.PhoneBook, contact models.Contact) error {
	return pb.db.Model(phonebook).Association("Contacts").Append(&contact)
}

func (pb *phoneBookGormRepository) DeleteContact(phonebook *models.PhoneBook, contactId uint) error {
	contact := &models.Contact{ID: contactId}
	return pb.db.Delete(contact).Error
}
func (pb *phoneBookGormRepository) GetContact(contactId uint) (*models.Contact, error) {
	var contact models.Contact
	err := pb.db.First(&contact, contactId).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (pb *phoneBookGormRepository) UpdateContact(phonebook *models.PhoneBook, contact *models.Contact) error {
	contact.PhoneBookID = phonebook.ID
	return pb.db.Save(contact).Error
}
