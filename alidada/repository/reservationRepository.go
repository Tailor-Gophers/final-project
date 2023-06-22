package repository

import (
	"alidada/db"
	"alidada/models"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	PlaceOrder(order *models.Order) error
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
	return nil //todo
}
