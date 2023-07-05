package services

import (
	"mockapi/models"
	"mockapi/repository"
	"time"
)

func NewFlightService(flightRepo repository.FlightRepository) FlightService {
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
	ReserveFlightCapacity(id int64) (*models.FlightClass, error)
	ReturnFlightCapacity(id int64) (*models.FlightClass, error)
	GetFlightByFilter(airline string, aircraft string, departure time.Time) ([]models.FlightClass, error)
	GetFlightBySort(order string) (*[]models.FlightClass, error)
	GetFlightClassByID(id int64) (*models.FlightClass, error)
}

type flightService struct {
	flightRepository repository.FlightRepository
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

func (s *flightService) ReserveFlightCapacity(id int64) (*models.FlightClass, error) {
	return s.flightRepository.ReserveFlightCapacity(id)
}

func (s *flightService) ReturnFlightCapacity(id int64) (*models.FlightClass, error) {
	return s.flightRepository.ReturnFlightCapacity(id)
}

func (s *flightService) GetFlightByFilter(airline string, aircraft string, departure time.Time) ([]models.FlightClass, error) {
	return s.flightRepository.GetFlightByFilter(airline, aircraft, departure)
}

func (s *flightService) GetFlightBySort(order string) (*[]models.FlightClass, error) {
	return s.flightRepository.GetFlightBySort(order)
}

func (s *flightService) GetFlightClassByID(id int64) (*models.FlightClass, error) {
	return s.flightRepository.GetFlightClassByID(id)
}
