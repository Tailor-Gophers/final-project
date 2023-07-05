package repository

import (
	"gorm.io/gorm"
	"qsms/db"
	"qsms/models"
)

type AdminRepository interface {
	AddNumber(number *models.Number) error
	SuspendUser(userId uint) error
	UnSuspendUser(userId uint) error
	CountUserMessages(userId uint) (int, error)
	GetAllMessages() ([]models.Message, error)
	GetMessageByID(messageId uint) (models.Message, error)
}

type adminGormRepository struct {
	db *gorm.DB
}

func NewGormAdminRepository() AdminRepository {
	return &adminGormRepository{
		db: db.GetDbConnection(),
	}
}

func (ar *adminGormRepository) AddNumber(number *models.Number) error {
	return ar.db.Create(number).Error
}

func (ar *adminGormRepository) SuspendUser(userId uint) error {
	return ar.db.Model(&models.User{}).Where("id = ?", userId).Update("Disable", true).Error
}

func (ar *adminGormRepository) UnSuspendUser(userId uint) error {
	return ar.db.Model(&models.User{}).Where("id = ?", userId).Update("Disable", false).Error
}

func (ar *adminGormRepository) CountUserMessages(userId uint) (int, error) {
	var messages []models.Message
	err := ar.db.Where("sender_id = ?", userId).Find(&messages).Error
	return len(messages), err
}

func (ar *adminGormRepository) GetAllMessages() ([]models.Message, error) {
	var messages []models.Message
	err := ar.db.Find(&messages).Error
	return messages, err
}

func (ar *adminGormRepository) GetMessageByID(messageId uint) (models.Message, error) {
	var message models.Message
	err := ar.db.First(&message, messageId).Error
	return message, err
}
