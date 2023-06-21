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
	GetCitiesList() ([]string, error)
	GetDaysList() ([]string, error)
	ReserveFlightCapacity(id int64, class string) (*models.FlightClass, error)
	ReturnFlightCapacity(id int64, class string) (*models.FlightClass, error)
	GetFlightByFilter(airline string, aircraft string, departure time.Time) ([]models.FlightClass, error)
	GetFlightBySort(order string) (*[]models.FlightClass, error)
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
	result := fl.db.Where("origin = ? and destination = ? and date(start_time) = date(?)", origin, destination, day).Find(&flights)
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

func (fl *flightGormRepository) GetCitiesList() ([]string, error) {
	var origin, destination []string
	result := fl.db.Model(&models.Flight{}).Distinct("origin").Pluck("origin", &origin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("planes not found")
	}
	result = fl.db.Model(&models.Flight{}).Distinct("destination").Pluck("destination", &destination)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("planes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	cities := append(origin, destination...)
	cities = removeDuplicateString(cities)

	return cities, nil
}

func (fl *flightGormRepository) GetDaysList() ([]string, error) {
	var startdate, enddate []string
	result := fl.db.Model(&models.Flight{}).Distinct("start_time").Pluck("start_time", &startdate)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("date not found")
	}
	result = fl.db.Model(&models.Flight{}).Distinct("end_time").Pluck("end_time", &enddate)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("date not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	dates := append(startdate, enddate...)
	dates = removeDuplicateString(dates)

	return dates, nil
}

func (fl *flightGormRepository) ReserveFlightCapacity(id int64, class string) (*models.FlightClass, error) {
	var flightClass *models.FlightClass
	err := fl.db.Where("id = ? and title = ?", id, class).Preload("Flight").First(&flightClass).Error
	if err != nil {
		return nil, err
	}

	if flightClass.Reserve != nil && *flightClass.Reserve == flightClass.Capacity {
		return nil, errors.New("flight capacity reached")
	} else {
		newReserve := *flightClass.Reserve + 1
		flightClass.Reserve = &newReserve
	}

	if err := fl.db.Save(flightClass).Error; err != nil {
		return nil, err
	}

	return flightClass, nil
}

func (fl *flightGormRepository) ReturnFlightCapacity(id int64, class string) (*models.FlightClass, error) {
	flightClass := &models.FlightClass{}
	err := fl.db.Where("flight_id = ? and title = ?", id, class).Preload("Flight").First(&flightClass).Error
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

func (fl *flightGormRepository) GetFlightByFilter(airline string, aircraft string, departure time.Time) ([]models.FlightClass, error) {
	var flightClass []models.FlightClass
	result := fl.db.
		Joins("join flights on flights.Id = flight_id").
		Where("flights.airline = ? and flights.aircraft = ? and date(flights.start_time) = date(?)", airline, aircraft, departure).
		Where("capacity-reserve != ?", 0).
		Preload("Flight").
		Find(&flightClass)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no flights found")
		}
		return nil, err
	}

	return flightClass, nil
}

func (fl *flightGormRepository) GetFlightBySort(order string) (*[]models.FlightClass, error) {
	var flights []models.FlightClass
	var result *gorm.DB

	if order == "decs" {
		result = fl.db.
			Joins("join flights on flights.Id = flight_classes.flight_id").
			Order("price desc, start_time, end_time-start_time").
			Preload("Flight").
			Find(&flights)
	} else {
		result = fl.db.
			Joins("join flights on flights.Id = flight_classes.flight_id").
			Order("price, start_time, end_time-start_time").
			Preload("Flight").
			Find(&flights)
	}

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no flights found")
		}
		return nil, err
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

func removeDuplicateString(strSlice []string) []string {
	// map to store unique keys
	keys := make(map[string]bool)
	returnSlice := []string{}
	for _, item := range strSlice {
		if _, value := keys[item]; !value {
			keys[item] = true
			returnSlice = append(returnSlice, item)
		}
	}
	return returnSlice
}
