package main

import (
	"alidada/controllers"
	"alidada/models"
	"alidada/repository"
	"alidada/services"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type sortTest struct {
	cancellationConditions       *[]models.CancellationCondition
	start                        int
	end                          int
	sortedCancellationConditions []models.CancellationCondition
}

var sortTests = []sortTest{
	sortTest{&[]models.CancellationCondition{
		{1, 10, "10 min", "", 90},
		{2, 20, "20 min", "", 80},
		{3, 20, "20 min", "", 70},
	}, 0, 2,
		[]models.CancellationCondition{
			{3, 20, "20 min", "", 70},
			{2, 20, "20 min", "", 80},
			{1, 10, "10 min", "", 90},
		},
	},
	sortTest{&[]models.CancellationCondition{
		{1, 10, "10 min", "", 10},
		{2, 20, "20 min", "", 20},
		{3, 20, "20 min", "", 30},
	}, 0, 2,
		[]models.CancellationCondition{
			{1, 10, "10 min", "", 10},
			{2, 20, "20 min", "", 20},
			{3, 20, "20 min", "", 30},
		},
	},
	sortTest{&[]models.CancellationCondition{
		{1, 10, "10 min", "", 80},
		{2, 20, "20 min", "", 10},
		{3, 20, "20 min", "", 12},
		{4, 20, "20 min", "", 70},
	}, 0, 3,
		[]models.CancellationCondition{
			{2, 20, "20 min", "", 10},
			{3, 20, "20 min", "", 12},
			{4, 20, "20 min", "", 70},
			{1, 10, "10 min", "", 80},
		},
	},
}

func TestSort(t *testing.T) {

	for j, test := range sortTests {
		output := services.Sort(test.cancellationConditions, test.start, test.end)
		for i, out := range output {
			if out.ID != test.sortedCancellationConditions[i].ID {
				t.Errorf("Output %d not equal to expected %d test number %d", out.ID, test.sortedCancellationConditions[i].ID, j+1)
			}
		}
	}
}

type penaltyTest struct {
	reservation                  *models.Reservation
	flightclass                  *models.FlightClass
	sortedCancellationConditions []models.CancellationCondition
	penalty                      int
	error                        error
}

var penaltyTests = []penaltyTest{
	penaltyTest{&models.Reservation{Price: 100},
		&models.FlightClass{Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 41)}},
		[]models.CancellationCondition{
			{2, 40, "40 min", "", 10},
			{3, 30, "30 min", "", 12},
			{4, 20, "20 min", "", 70},
			{1, 10, "10 min", "", 80},
		},
		10, nil,
	},

	penaltyTest{&models.Reservation{Price: 200},
		&models.FlightClass{Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 34)}},
		[]models.CancellationCondition{
			{2, 40, "40 min", "", 10},
			{3, 30, "30 min", "", 12},
			{4, 20, "20 min", "", 70},
			{1, 10, "10 min", "", 80},
		},
		24, nil,
	},

	penaltyTest{&models.Reservation{Price: 300},
		&models.FlightClass{Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 24)}},
		[]models.CancellationCondition{
			{2, 40, "40 min", "", 10},
			{3, 30, "30 min", "", 12},
			{4, 20, "20 min", "", 70},
			{1, 10, "10 min", "", 80},
		},
		210, nil,
	},
	penaltyTest{&models.Reservation{Price: 400},
		&models.FlightClass{Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 12)}},
		[]models.CancellationCondition{
			{2, 40, "40 min", "", 10},
			{3, 30, "30 min", "", 12},
			{4, 20, "20 min", "", 70},
			{1, 10, "10 min", "", 80},
		},
		320, nil,
	},
	penaltyTest{&models.Reservation{Price: 500},
		&models.FlightClass{Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 8)}},
		[]models.CancellationCondition{
			{2, 40, "40 min", "", 10},
			{3, 30, "30 min", "", 12},
			{4, 20, "20 min", "", 70},
			{1, 10, "10 min", "", 80},
		},
		100, errors.New("None of the cancellation conditions are available for you"),
	},
}

func TestPenalty(t *testing.T) {
	for j, test := range penaltyTests {
		penalty, error := services.PenaltyCalculation(test.reservation, test.flightclass, test.sortedCancellationConditions)
		if penalty != test.penalty {
			t.Errorf("Output %d not equal to expected %d test number %d", penalty, test.penalty, j+1)
		}
		if (error != nil && test.error == nil) || (error == nil && test.error != nil) {
			t.Errorf("Error not equal test number %d", j+1)

		}

	}
}

type generatePdfTest struct {
	reservations []models.Reservation
	saveTo       string
	error        error
}

var generatePdfTests = []generatePdfTest{
	generatePdfTest{
		[]models.Reservation{
			{Model: models.Model{ID: 1}, FlightClassID: 1, OrderID: 1, Price: 1, IsCancelled: false,
				Passenger: models.Passenger{FirstName: "eghbal", LastName: "shirasb", DateOfBirth: "26/9/77",
					Nationality: "irani", PassportNumber: "1916784171"}, Confirmed: true,
				FlightClass: models.FlightClass{Model: models.Model{ID: 1}, Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 41), Origin: "tehran", Destination: "mashhad"}, Title: "ecconomy"}},
		},
		"pdf/test_1.pdf",
		nil,
	},
	generatePdfTest{
		[]models.Reservation{
			{Model: models.Model{ID: 1}, FlightClassID: 1, OrderID: 1, Price: 1, IsCancelled: false,
				Passenger: models.Passenger{FirstName: "eghbal", LastName: "shirasb", DateOfBirth: "26/9/77",
					Nationality: "irani", PassportNumber: "1916784171"}, Confirmed: true,
				FlightClass: models.FlightClass{Model: models.Model{ID: 1}, Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 41), Origin: "tehran", Destination: "mashhad"}, Title: "ecconomy"}},
			{Model: models.Model{ID: 1}, FlightClassID: 1, OrderID: 1, Price: 1, IsCancelled: false,
				Passenger: models.Passenger{FirstName: "armin", LastName: "armin", DateOfBirth: "26/9/77",
					Nationality: "irani", PassportNumber: "1916784171"}, Confirmed: true,
				FlightClass: models.FlightClass{Model: models.Model{ID: 1}, Flight: &models.Flight{StartTime: time.Now().Add(-1 * time.Minute * 41), Origin: "tehran", Destination: "mashhad"}, Title: "ecconomy"}},
		},
		"pdf/test_2.pdf",
		nil,
	},
}

func TestGeneratePdf(t *testing.T) {
	for j, test := range generatePdfTests {
		error := services.GeneratePdf(test.reservations, test.saveTo)
		if (error != nil && test.error == nil) || (error == nil && test.error != nil) {
			t.Errorf("Error not equal test number %d", j+1)
		}
		if error != nil {
			if _, err := os.Stat(test.saveTo); errors.Is(err, os.ErrNotExist) {
				t.Errorf("Error pdf not generated test number %d", j+1)
			}
		}
	}

}

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
				ID: 1,
			},
			Title:    "Class-A",
			Price:    1900,
			Capacity: 50,
			FlightId: 6,
		},
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

type mockUserRepository struct {
	mock.Mock
	repository.UserRepository
}

func (m *mockUserRepository) CreatePassenger(passenger *models.Passenger) error {
	args := m.Called(passenger)
	return args.Error(0)
}

func (m *mockUserRepository) GetPassengers(user *models.User) ([]models.Passenger, error) {
	args := m.Called(user)
	return args.Get(0).([]models.Passenger), args.Error(1)
}

func TestCreatePassenger(t *testing.T) {
	userRepo := &mockUserRepository{}
	service := services.UserServicet{UserRepository: userRepo}

	mockPassenger := &models.Passenger{
		Model: models.Model{
			ID: 2,
		},
		UserID:         2,
		FirstName:      "Behzad",
		LastName:       "Ebram",
		Gender:         false,
		DateOfBirth:    "000-1",
		Nationality:    "Iranian",
		PassportNumber: "123456712",
	}

	userRepo.On("CreatePassenger", mockPassenger).Return(nil)

	err := service.CreatePassenger(mockPassenger)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestGetPassengers(t *testing.T) {
	userRepo := &mockUserRepository{}
	service := services.UserServicet{UserRepository: userRepo}

	mockUser := &models.User{
		Model: models.Model{
			ID: 2,
		},
		Email:       "amir@me.com",
		UserName:    "amir",
		Password:    "amir1234",
		FirstName:   "Amir",
		LastName:    "Val",
		PhoneNumber: "1234567890",
		IsAdmin:     false,
		Passengers:  []models.Passenger{},
		Tokens:      []models.Token{},
	}

	userRepo.On("GetPassengers", mockUser).Return([]models.Passenger{}, nil)

	_, err := service.GetPassengers(mockUser)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}
