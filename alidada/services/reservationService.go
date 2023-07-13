package services

import (
	"alidada/utils"

	"alidada/models"
	"alidada/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ReservationService interface {
	Reserve(passengers []uint, flightClassId uint) (*models.Order, error)
	SetAuthorityPair(authority string, userId uint) error
	GetOrderByAuthority(authority string) (*models.Order, error)
	ConfirmOrder(orderId uint, refId int) error
}

type reservationService struct {
	reservationRepository repository.ReservationRepository
}

func NewReservationService() ReservationService {
	return &reservationService{
		reservationRepository: repository.NewGormReservationRepository(),
	}
}

func (rs *reservationService) Reserve(passengers []uint, flightClassId uint) (*models.Order, error) {
	var flightClass models.FlightClass

	url := fmt.Sprintf("%s/flight_class/%d", utils.ENV("MOCK_URL"), flightClassId)
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get flight with id %d from mock api", flightClassId))
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&flightClass)

	flightClass.Price = GetPrice(flightClass)
	if int(flightClass.Capacity-*flightClass.Reserve) < len(passengers) {
		return nil, errors.New(fmt.Sprintf("Not enouph capacity in flight %d ", flightClassId))
	}
	err = rs.reservationRepository.CanReserve(passengers, flightClassId)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		Reservations: []models.Reservation{},
		Price:        0,
		OrderTime:    time.Now(),
		Confirmed:    false,
	}

	for _, passenger := range passengers {

		reservation := models.Reservation{
			PassengerID:   passenger,
			FlightClassID: flightClassId,
			Price:         flightClass.Price,
			IsCancelled:   false,
		}
		order.Reservations = append(order.Reservations, reservation)
		order.Price += reservation.Price
	}

	err = rs.reservationRepository.PlaceOrder(order)
	if err != nil {
		return nil, errors.New("failed to place order")
	}
	return order, err
}

func (rs *reservationService) SetAuthorityPair(authority string, orderId uint) error {
	return rs.reservationRepository.SetAuthorityPair(authority, orderId)
}

func (rs *reservationService) GetOrderByAuthority(authority string) (*models.Order, error) {
	return rs.reservationRepository.GetOrderByAuthority(authority)
}

func (rs *reservationService) ConfirmOrder(orderId uint, refId int) error {
	return rs.reservationRepository.ConfirmOrder(orderId, refId)
}

func SmartPrice(flightclasses []models.FlightClass) []models.FlightClass {
	for i, flightclass := range flightclasses {
		flightclasses[i].Price = GetPrice(flightclass)
	}
	return flightclasses
}

func GetPrice(flightclass models.FlightClass) uint {

	if int(*flightclass.Reserve)*100/int(flightclass.Capacity) < 10 {
		return uint(float32(flightclass.Price) * 0.8)
	}
	if int(*flightclass.Reserve)*100/int(flightclass.Capacity) < 20 {
		return uint(float32(flightclass.Price) * 0.9)
	}
	if int(*flightclass.Reserve)*100/int(flightclass.Capacity) < 50 {
		return uint(float32(flightclass.Price) * 0.95)
	}
	if *flightclass.Reserve*100/flightclass.Capacity < 60 {
		return uint(float32(flightclass.Price) * 1.1)
	}
	if *flightclass.Reserve*100/flightclass.Capacity < 70 {
		return uint(float32(flightclass.Price) * 1.3)
	}
	if *flightclass.Reserve*100/flightclass.Capacity < 80 {
		return uint(float32(flightclass.Price) * 1.5)
	}
	return uint(float32(flightclass.Price) * 2)
}
