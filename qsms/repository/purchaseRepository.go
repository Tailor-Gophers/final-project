package repository

import (
	"gorm.io/gorm"
	"qsms/db"
	"qsms/models"
	"time"
)

type PurchaseRepository interface {
	GetNumberByID(numberId uint) (models.Number, error)
	SetUserBalance(user *models.User, balance int) error
	UpdateNumber(user *models.User, numberId uint) error
	RestoreNumber(numberId uint) error
	PlaceRent(rent *models.Rent) error
	DropRent(rentId uint) error
	GetAllRents() ([]models.Rent, error)
	GetRentByIDs(userId, numberId uint) (models.Rent, error)
	GetUserById(userId uint) (*models.User, error)
	UpdateRentDate(rentId uint, time time.Time) error
	GetRentByID(rentId uint) (*models.Rent, error)
	UpdateUserMainNumber(userId uint) error
}

type purchaseGormRepository struct {
	db *gorm.DB
}

func NewGormPurchaseRepository() PurchaseRepository {
	return &purchaseGormRepository{
		db: db.GetDbConnection(),
	}
}

func (pr *purchaseGormRepository) GetNumberByID(numberId uint) (models.Number, error) {
	var number models.Number
	err := pr.db.First(&number, numberId).Error
	return number, err
}

func (pr *purchaseGormRepository) SetUserBalance(user *models.User, balance int) error {
	return pr.db.Model(user).Update("Balance", balance).Error
}

func (pr *purchaseGormRepository) UpdateNumber(user *models.User, numberId uint) error {
	return pr.db.Model(&models.Number{}).Where("id = ?", numberId).
		Updates(map[string]interface{}{"UserID": user.ID, "Active": true}).Error
}

func (pr *purchaseGormRepository) RestoreNumber(numberId uint) error {
	return pr.db.Model(&models.Number{}).Where("id = ?", numberId).
		Updates(map[string]interface{}{"UserID": 0, "Active": false}).Error
}

func (pr *purchaseGormRepository) PlaceRent(rent *models.Rent) error {
	return pr.db.Create(rent).Error
}

func (pr *purchaseGormRepository) GetAllRents() ([]models.Rent, error) {
	var rents []models.Rent
	err := pr.db.Find(&rents).Error
	return rents, err
}

func (pr *purchaseGormRepository) DropRent(rentId uint) error {
	return pr.db.Delete(&models.Rent{}, rentId).Error
}

func (pr *purchaseGormRepository) GetRentByID(rentId uint) (*models.Rent, error) {
	var rent models.Rent
	err := pr.db.First(&rent, rentId).Error
	return &rent, err
}

func (pr *purchaseGormRepository) GetRentByIDs(userId, numberId uint) (models.Rent, error) {
	var rent models.Rent
	err := pr.db.Where(map[string]interface{}{"user_id": userId, "number_id": numberId}).First(&rent).Error
	return rent, err
}

func (pr *purchaseGormRepository) GetUserById(userId uint) (*models.User, error) {
	var user *models.User
	err := pr.db.First(user, userId).Error
	return user, err
}

func (pr *purchaseGormRepository) UpdateRentDate(rentId uint, time time.Time) error {
	return pr.db.Model(&models.Rent{}).Where("id = ?", rentId).Update("LastPaid", time).Error
}

func (pr *purchaseGormRepository) UpdateUserMainNumber(userId uint) error {
	return pr.db.Model(&models.User{}).Where("id = ?", userId).Update("MainNumberID", 0).Error
}
