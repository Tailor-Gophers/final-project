package controllers_test

import (
	"alidada/controllers"
	"alidada/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type FlightClassComparison struct {
	Model    models.Model
	Title    string
	Price    uint
	Capacity uint
	FlightId uint
}

func TestSearchFlightsDay(t *testing.T) {
	e := echo.New()

	startTime1, _ := time.Parse(time.RFC3339, "2020-02-08T15:49:41+03:30")
	endTime1, _ := time.Parse(time.RFC3339, "2020-02-09T18:49:41+03:30")
	expectedFlights := &[]models.Flight{
		{
			Model: models.Model{
				ID: 2,
			},
			Origin:      "Shiraz",
			Destination: "Tehran",
			StartTime:   startTime1,
			EndTime:     endTime1,
			Airline:     "homa",
			Aircraft:    "Boeing737",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/flights/search?origin=Shiraz&destination=Tehran&date=2020-02-08", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	fc := &controllers.FlightController{}
	err := fc.SearchFlightsDay(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var flights []models.Flight
	err = json.Unmarshal(rec.Body.Bytes(), &flights)
	assert.NoError(t, err)

	if !reflect.DeepEqual(*expectedFlights, flights) {
		t.Errorf("Expected flights: %v, but got: %v", *expectedFlights, flights)
	}
}

func TestSearchFlightsSort(t *testing.T) {
	e := echo.New()

	expectedFlightClass := []FlightClassComparison{
		{
			Model: models.Model{
				ID: 4,
			},
			Title:    "Class-A",
			Price:    1900,
			Capacity: 50,
			FlightId: 2,
		},
		{
			Model: models.Model{
				ID: 2,
			},
			Title:    "Class-B",
			Price:    1700,
			Capacity: 50,
			FlightId: 6,
		},
		{
			Model: models.Model{
				ID: 1,
			},
			Title:    "Class-A",
			Price:    1300,
			Capacity: 50,
			FlightId: 6,
		},
		{
			Model: models.Model{
				ID: 3,
			},
			Title:    "Class-C",
			Price:    1300,
			Capacity: 50,
			FlightId: 6,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/flights/sort/price?order=desc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("orderBy")
	c.SetParamValues("price")

	fc := &controllers.FlightController{}
	err := fc.SearchFlightsSort(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var flightclass []models.FlightClass
	err = json.Unmarshal(rec.Body.Bytes(), &flightclass)
	assert.NoError(t, err)

	var flightclassComparison []FlightClassComparison
	for _, fc := range flightclass {
		flightclassComparison = append(flightclassComparison, FlightClassComparison{
			Model:    fc.Model,
			Title:    fc.Title,
			Price:    fc.Price,
			Capacity: fc.Capacity,
			FlightId: fc.Flight.ID,
		})
	}

	if !reflect.DeepEqual(expectedFlightClass, flightclassComparison) {
		t.Errorf("Expected flights: %v, but got: %v", expectedFlightClass, flightclassComparison)
	}
}

func TestFilterFlights(t *testing.T) {
	e := echo.New()

	expectedFlightClass := []FlightClassComparison{
		{
			Model: models.Model{
				ID: 4,
			},
			Title:    "Class-A",
			Price:    1900,
			Capacity: 50,
			FlightId: 2,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/flights/filter?airline=homa&aircraft=Boeing737&departure=2020-02-08&capacity=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("filterBy")
	c.SetParamValues("price")

	fc := &controllers.FlightController{}
	err := fc.FilterFlights(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var flightclass []models.FlightClass
	err = json.Unmarshal(rec.Body.Bytes(), &flightclass)
	assert.NoError(t, err)

	var flightclassComparison []FlightClassComparison
	for _, fc := range flightclass {
		flightclassComparison = append(flightclassComparison, FlightClassComparison{
			Model:    fc.Model,
			Title:    fc.Title,
			Price:    fc.Price,
			Capacity: fc.Capacity,
			FlightId: fc.Flight.ID,
		})
	}

	if !reflect.DeepEqual(expectedFlightClass, flightclassComparison) {
		t.Errorf("Expected flights: %v, but got: %v", expectedFlightClass, flightclassComparison)
	}
}
