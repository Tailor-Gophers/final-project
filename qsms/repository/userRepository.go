package repository

import (
	"errors"
	"fmt"
	"qsms/db"
	"qsms/models"
	"qsms/utils"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUserName(username string) (*models.User, error)
	GetUserByEmail(username string) (*models.User, error)
	GetUserById(userId uint) (*models.User, error)
	DeleteUser(userId uint) error
	SaveToken(user *models.User, token string) error
	UserByToken(token string) (*models.User, error)
	LogOut(token string) error
	UpdateBalance(userId uint, amount int) error
	AddContact(contact *models.Contact) error
	DeleteContact(user *models.User, contactId uint) error
	GetContact(contactId uint) (*models.Contact, error)
	UpdateContact(user *models.User, contact *models.Contact) error
	GetUserByID(userId uint) (*models.User, error)
	GetAvailablePhoneNumbers() ([]models.Number, error)
	SetMainNumber(user *models.User, numberId uint) error
	CreateTemplate(template *models.Template) error
	GetTemplate(templateId uint) (*models.Template, error)
}

type userGormRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() UserRepository {
	return &userGormRepository{
		db: db.GetDbConnection(),
	}
}

func (ur *userGormRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userGormRepository) GetUserByUserName(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userGormRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userGormRepository) DeleteUser(userId uint) error {
	return ur.db.Delete(&models.User{}, userId).Error
}

func (ur *userGormRepository) GetUserById(userId uint) (*models.User, error) {
	var user models.User
	result := ur.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *userGormRepository) SaveToken(user *models.User, token string) error {

	hashed, err := utils.HashToken(token)
	if err != nil {
		return err
	}
	AccessToken := models.Token{UserId: user.ID, Token: hashed, ExpiresAt: time.Now().Add(time.Hour * 24)}

	result := ur.db.Create(&AccessToken)

	return result.Error
}

func (ur *userGormRepository) UserByToken(token string) (*models.User, error) {
	var AccessToken models.Token
	var User models.User

	hashed, err := utils.HashToken(token)
	if err != nil {
		return nil, err
	}
	err = ur.db.Where("token = ?", hashed).Where("expires_at > ?", time.Now()).First(&AccessToken).Error
	if err != nil {
		return nil, err
	}
	err = ur.db.Preload("Numbers").Preload("PhoneBooks").Preload("PhoneBooks.Numbers").Preload("Contacts").Preload("Templates").
		Where("id = ?", AccessToken.UserId).First(&User).Error
	if err != nil {
		return nil, err
	}
	return &User, nil
}

func (ur *userGormRepository) LogOut(token string) error {
	var AccessToken models.Token

	hashed, err := utils.HashToken(token)
	if err != nil {
		return err
	}
	err = ur.db.Where("token = ?", hashed).Where("expires_at > ?", time.Now()).First(&AccessToken).Error

	ur.db.Where("token = ?", hashed).Where("expires_at > ?", time.Now()).Delete(&AccessToken)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userGormRepository) UpdateBalance(userId uint, amount int) error {
	return ur.db.Model(&models.User{}).Where("id = ?", userId).Update("balance", amount).Error
}

func (ur *userGormRepository) AddContact(contact *models.Contact) error {
	return ur.db.Create(contact).Error
}

func (ur *userGormRepository) CreateTemplate(template *models.Template) error {
	return ur.db.Create(template).Error
}

func (ur *userGormRepository) DeleteContact(user *models.User, contactId uint) error {
	return ur.db.Delete(&models.Contact{}, contactId).Error
}

func (ur *userGormRepository) GetContact(contactId uint) (*models.Contact, error) {
	var contact models.Contact
	err := ur.db.First(&contact, contactId).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (ur *userGormRepository) UpdateContact(user *models.User, contact *models.Contact) error {
	contact.UserID = user.ID
	return ur.db.Save(contact).Error
}

func (ur *userGormRepository) GetUserByID(userId uint) (*models.User, error) {
	var user models.User
	err := ur.db.Preload("Numbers").Preload("PhoneBooks").Preload("PhoneBooks.Numbers").Preload("Contacts").Preload("Templates").First(&user, userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userGormRepository) GetAvailablePhoneNumbers() ([]models.Number, error) {
	var numbers []models.Number
	err := ur.db.Where("active = ?", 0).Find(&numbers).Error
	return numbers, err
}

func (ur *userGormRepository) SetMainNumber(user *models.User, numberId uint) error {
	return ur.db.Model(&user).Update("MainNumberID", numberId).Error
}

func (ur *userGormRepository) GetTemplate(templateId uint) (*models.Template, error) {
	var template models.Template
	err := ur.db.Where("id = ?", templateId).First(&template).Error
	return &template, err
}
