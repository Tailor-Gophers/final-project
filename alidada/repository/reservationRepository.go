package repository

import (
	"alidada/db"
	"alidada/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	PlaceOrder(order *models.Order) error
	SetAuthorityPair(authority string, id uint) error
	GetOrderByAuthority(authority string) (*models.Order, error)
	ConfirmOrder(order uint, refId int) error
}
type reservationGormRepository struct {
	db *gorm.DB
}

func NewGormReservationRepository() ReservationRepository {
	return &reservationGormRepository{
		db: db.GetDbConnection(),
	}
}

func (rr *reservationGormRepository) PlaceOrder(order *models.Order) error {
	return rr.db.Create(order).Error
}

func (rr *reservationGormRepository) SetAuthorityPair(authority string, id uint) error {
	pair := &models.AuthorityPair{
		OrderID:   id,
		Authority: authority,
	}
	return rr.db.Create(pair).Error
}

func (rr *reservationGormRepository) GetOrderByAuthority(authority string) (*models.Order, error) {
	var order *models.Order
	err := rr.db.Preload(clause.Associations).Where("authority = ?", authority).First(order).Error

	return order, err
}

func (rr *reservationGormRepository) ConfirmOrder(orderId uint, refId int) error {
	return rr.db.Model(&models.Order{}).Where("id = ", orderId).Updates(map[string]interface{}{"ref_id": refId, "confirmed": true}).Error
}
