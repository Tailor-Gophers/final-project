package repository

import (
	"errors"
	"fmt"
	"mockapi/db"
	"mockapi/models"
	"time"

	"gorm.io/gorm"
)

type FlightRepository interface {
	GetFlightsByCityAndDate(origin string, destination string, day time.Time) ([]models.Flight, error)
	GetFlightByID(id int64) (*models.Flight, error)
	GetPlanesList() ([]string, error)
	GetCitiesList() (*[]models.Flight, *[]models.Flight, error)
	GetDaysList() ([]models.Flight, error)
	ReserveFlightCapacity(id int64, count int) (*models.FlightClass, error)
	ReturnFlightCapacity(id int64) (*models.FlightClass, error)
	GetFlightByFilter(airline string, aircraft string, departure time.Time, capacity uint) ([]models.FlightClass, error)
	GetFlightBySort(orderBy string, order string) (*[]models.FlightClass, error)
	GetFlightClassByID(id int64) (*models.FlightClass, error)
}

type flightGormRepository struct {
	db *gorm.DB
}

func NewGormFlightRepository() FlightRepository {
	return &flightGormRepository{
		db: db.GetDbConnection(),
	}
}

func (fl *flightGormRepository) GetFlightsByCityAndDate(origin string, destination string, day time.Time) ([]models.Flight, error) {
	var flights []models.Flight
	result := fl.db.Where("origin = ? and destination = ? and date(start_time) = date(?)", origin, destination, day).Order("id").Find(&flights)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("flights not found")
		}
		return nil, err
	}
	return flights, nil
}

func (fl *flightGormRepository) GetFlightByID(id int64) (*models.Flight, error) {
	var flight models.Flight
	result := fl.db.First(&flight, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("flight not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &flight, nil
}

func (fl *flightGormRepository) GetPlanesList() ([]string, error) {
	var planes []string
	result := fl.db.Model(&models.Flight{}).Distinct("aircraft").Pluck("aircraft", &planes)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("planes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return planes, nil
}

func (fl *flightGormRepository) GetCitiesList() (*[]models.Flight, *[]models.Flight, error) {
	var origin, destination *[]models.Flight
	result := fl.db.Select("id", "origin", "start_time", "end_time", "airline", "aircraft").Find(&origin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil, fmt.Errorf("origins not found")
	}

	result = fl.db.Select("id", "destination", "start_time", "end_time", "airline", "aircraft").Find(&destination)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil, fmt.Errorf("cities not found")
	}

	if result.Error != nil {
		return nil, nil, result.Error
	}

	return origin, destination, nil
}

func (fl *flightGormRepository) GetDaysList() ([]models.Flight, error) {
	var days []models.Flight
	result := fl.db.Select("id", "origin", "destination", "start_time", "end_time", "airline", "aircraft").Find(&days)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("days not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return days, nil
}

func (fl *flightGormRepository) ReserveFlightCapacity(id int64, count int) (*models.FlightClass, error) {
	var flightClass *models.FlightClass
	err := fl.db.Where("id = ?", id).Preload("Flight").First(&flightClass).Error
	if err != nil {
		return nil, err
	}

	if flightClass.Reserve != nil && *flightClass.Reserve+uint(count) > flightClass.Capacity {
		return nil, errors.New("flight capacity reached")
	} else {
		newReserve := *flightClass.Reserve + uint(count)
		flightClass.Reserve = &newReserve
	}

	if err := fl.db.Save(flightClass).Error; err != nil {
		return nil, err
	}

	return flightClass, nil
}

func (fl *flightGormRepository) ReturnFlightCapacity(id int64) (*models.FlightClass, error) {
	flightClass := &models.FlightClass{}
	err := fl.db.Where("flight_id = ?", id).Preload("Flight").First(&flightClass).Error
	if err != nil {
		return nil, err
	}

	if *(flightClass.Reserve) == 0 {
		return nil, errors.New("flight capacity is already empty")
	} else {
		*flightClass.Reserve--
	}

	if err := fl.db.Save(flightClass).Error; err != nil {
		return nil, err
	}

	return flightClass, nil
}

func (fl *flightGormRepository) GetFlightByFilter(airline string, aircraft string, departure time.Time, capacity uint) ([]models.FlightClass, error) {
	var flightClass []models.FlightClass

	query := fl.db.Joins("join flights on flights.Id = flight_id")

	if airline != "" {
		query = query.Where("flights.airline = ?", airline)
	}
	if aircraft != "" {
		query = query.Where("flights.aircraft = ?", aircraft)
	}
	if !departure.IsZero() {
		query = query.Where("date(flights.start_time) = date(?)", departure)
	}
	if capacity != 0 {
		query = query.Where("capacity-reserve >= ?", capacity)
	}

	result := query.Preload("Flight").Find(&flightClass)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no flights found")
		}
		return nil, err
	}

	return flightClass, nil
}

func (fl *flightGormRepository) GetFlightBySort(orderBy string, order string) (*[]models.FlightClass, error) {
	var flights []models.FlightClass

	if orderBy == "" {
		order = "asc"
	}

	result := fl.db.Joins("join flights on flights.Id = flight_id").
		Order(orderBy + " " + order).
		Preload("Flight").
		Find(&flights)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("no flights found")
	}
	return &flights, nil
}

func (fl *flightGormRepository) GetFlightPrice(id int64) (models.FlightClass, error) {
	var flightC models.FlightClass
	err := fl.db.Where("flight_id = ?", id).First(&flightC).Error
	if err != nil {
		return models.FlightClass{}, err
	}
	return flightC, nil
}

func (fl *flightGormRepository) GetFlightClassByID(id int64) (*models.FlightClass, error) {
	var flightClass models.FlightClass
	result := fl.db.Preload("Flight").First(&flightClass, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("flight not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &flightClass, nil
}
