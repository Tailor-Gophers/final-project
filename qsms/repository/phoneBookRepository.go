package repository

import (
	"qsms/db"
	"qsms/models"

	"gorm.io/gorm"
)

type PhoneBookRepository interface {
	CreatePhoneBook(user *models.User, phonebook *models.PhoneBook) error
	GetPhoneBook(phonebookId uint) (*models.PhoneBook, error)
	UpdatePhoneBook(phonebook *models.PhoneBook) error
	DeletePhoneBook(phonebookId uint) error
}

type phoneBookGormRepository struct {
	db *gorm.DB
}

func NewGormPhoneBookRepository() PhoneBookRepository {
	return &phoneBookGormRepository{
		db: db.GetDbConnection(),
	}
}

func (pb *phoneBookGormRepository) CreatePhoneBook(user *models.User, phonebook *models.PhoneBook) error {
	if err := pb.db.Preload("PhoneBooks").First(user, user.ID).Error; err != nil {
		return err
	}
	phonebook.UserID = user.ID
	if err := pb.db.Create(phonebook).Error; err != nil {
		return err
	}
	return nil
}

func (pb *phoneBookGormRepository) GetPhoneBook(phonebookId uint) (*models.PhoneBook, error) {
	var phonebook models.PhoneBook
	err := pb.db.First(&phonebook, phonebookId).Error
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
