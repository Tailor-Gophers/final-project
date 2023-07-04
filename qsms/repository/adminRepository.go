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
