package main

import "testing"
import "alidada/services"
import "alidada/models"
import "time"
import "errors"
import "fmt"
import "os"

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
			fmt.Println(error)
			t.Errorf("Error not equal test number %d", j+1)
		}
		if error != nil {
			if _, err := os.Stat(test.saveTo); errors.Is(err, os.ErrNotExist) {
				t.Errorf("Error pdf not generated test number %d", j+1)
			}
		}
	}

}
