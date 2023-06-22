package services

import (
	"alidada/models"
	"alidada/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ReservationService interface {
	Reserve(passengers []uint, flightClassId uint) error
}

type reservationService struct {
	reservationRepository repository.ReservationRepository
}

func NewReservationService() ReservationService {
	return &reservationService{
		reservationRepository: repository.NewGormReservationRepository(),
	}
}

func (rs *reservationService) Reserve(passengers []uint, flightClassId uint) error {
	var flightClass models.FlightClass

	url := fmt.Sprintf("http://localhost:3001/flight_class/%d", flightClassId)
	res, err := http.Get(url)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get flight with id %d from mock api", flightClassId))
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&flightClass)

	if int(flightClass.Capacity-*flightClass.Reserve) < len(passengers) {
		return errors.New(fmt.Sprintf("Not enouph capacity in flight %d ", flightClassId))
	}

	order := &models.Order{
		Reservations: []models.Reservation{},
		Price:        0,
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
		return errors.New("failed to place order")
	}
	return err
}
