package flightService

import (
	"mockapi/models"
	"mockapi/repository/flightRepository"
	"time"
)

func NewFlightService(flightRepo flightRepository.FlightRepository) FlightService {
	return &flightService{
		flightRepository: flightRepo,
	}
}

type FlightService interface {
	GetFlight(id int64) (*models.Flight, error)
	GetFlightByDate(origin string, destination string, day time.Time) ([]models.Flight, error)
	GetPlanesList() ([]string, error)
	GetCitiesList() ([]string, error)
	GetDaysList() ([]string, error)
	ReserveFlightCapacity(id int64) (*models.Flight, error)
	ReturnFlightCapacity(id int64) (*models.Flight, error)
}

type flightService struct {
	flightRepository flightRepository.FlightRepository
}

func (s *flightService) GetFlight(id int64) (*models.Flight, error) {
	return s.flightRepository.GetFlightByID(id)
}

func (s *flightService) GetFlightByDate(origin string, destination string, day time.Time) ([]models.Flight, error) {
	return s.flightRepository.GetFlightsByCityAndDate(origin, destination, day)
}

func (s *flightService) GetPlanesList() ([]string, error) {
	return s.flightRepository.GetPlanesList()
}

func (s *flightService) GetCitiesList() ([]string, error) {
	return s.flightRepository.GetCitiesList()
}

func (s *flightService) GetDaysList() ([]string, error) {
	return s.flightRepository.GetDaysList()
}

func (s *flightService) ReserveFlightCapacity(id int64) (*models.Flight, error) {
	return s.flightRepository.ReserveFlightCapacity(id)
}

func (s *flightService) ReturnFlightCapacity(id int64) (*models.Flight, error) {
	return s.flightRepository.ReturnFlightCapacity(id)
}
