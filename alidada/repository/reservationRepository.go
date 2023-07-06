package repository

import (
	"alidada/db"
	"alidada/models"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	PlaceOrder(order *models.Order) error
	SetAuthorityPair(authority string, id uint) error
	GetOrderByAuthority(authority string) (*models.Order, error)
	ConfirmOrder(order uint, refId int) error
	CanReserve(passengers []uint, flightClassId uint) error
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

func (rr *reservationGormRepository) CanReserve(passengers []uint, flightClassId uint) error {
	var reservation models.Reservation
	err := rr.db.Where("passenger_id IN ?", passengers).Where("flight_class_id = ?", flightClassId).Where("confirmed =?", true).First(&reservation).Error
	fmt.Println(reservation)
	if err == nil {
		return errors.New(fmt.Sprintf("reservation alredy exist : id:%d", reservation.PassengerID))
	} else {
		return nil
	}
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
	err := rr.db.Joins("JOIN authority_pairs ON authority_pairs.order_id = orders.id ").Where("authority_pairs.authority = ?", authority).First(&order).Error
	return order, err
}

func (rr *reservationGormRepository) ConfirmOrder(orderId uint, refId int) error {

	err := rr.db.Model(&models.Reservation{}).Where("order_id = ?", orderId).Updates(map[string]interface{}{"confirmed": true}).Error
	if err != nil {
		return err
	}
	err = rr.db.Model(&models.Order{}).Where("id = ?", orderId).Updates(map[string]interface{}{"ref_id": refId, "confirmed": true}).Error
	if err != nil {
		return err
	}
	var reservations []models.Reservation
	rr.db.Where("order_id = ?", orderId).Find(&reservations)

	url := fmt.Sprintf("http://localhost:3001/flights/%d/reserve/%d", reservations[0].FlightClassID, len(reservations))
	res, err := http.Post(url, "", nil)
	if err != nil {
		return fmt.Errorf("Failed to reserve flights from mockapi")
	}
	defer res.Body.Close()
	return nil

}
