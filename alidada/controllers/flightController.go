package controllers

import (
	"alidada/models"
	"alidada/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	echo "github.com/labstack/echo/v4"
)

type FlightController struct{}

func NewFlightController() *FlightController {
	return &FlightController{}
}

func (f *FlightController) SearchFlightsDay(c echo.Context) error {
	origin := c.QueryParam("origin")
	destination := c.QueryParam("destination")
	dateStr := c.QueryParam("date")

	url := fmt.Sprintf("%s/flights/%s/%s/%s", utils.ENV("MOCK_URL"), origin, destination, dateStr)
	res, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get flights from mockapi",
		})
	}
	defer res.Body.Close()

	var flights []models.Flight
	err = json.NewDecoder(res.Body).Decode(&flights)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to decode flights from mockapi response",
		})
	}

	return c.JSON(http.StatusOK, flights)
}

func (f *FlightController) SearchFlightsSort(c echo.Context) error {
	orderBy := c.Param("orderBy")
	order := c.QueryParam("order")
	smart := c.QueryParam("smart")

	url := fmt.Sprintf("%s/flights/sort/%s?order=%s&smart=%s", utils.ENV("MOCK_URL"), orderBy, order, smart)
	res, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get flights from mockapi",
		})
	}
	defer res.Body.Close()

	var flightclass []models.FlightClass
	err = json.NewDecoder(res.Body).Decode(&flightclass)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to decode flights from mockapi response",
		})
	}

	if smart == "1" {
		fakeFlightClass := make([]models.FlightClass, len(flightclass))
		copy(fakeFlightClass, flightclass)

		sort.SliceStable(fakeFlightClass, func(i, j int) bool {
			return *fakeFlightClass[i].Reserve < *fakeFlightClass[j].Reserve
		})

		toRemove := []models.FlightClass{fakeFlightClass[0], fakeFlightClass[1]}

		result := []models.FlightClass{}
		for _, fc := range flightclass {
			found := false
			for _, tr := range toRemove {
				if fc == tr {
					found = true
					break
				}
			}
			if !found {
				result = append(result, fc)
			}
		}

		flightclass = append([]models.FlightClass{fakeFlightClass[0], fakeFlightClass[1]}, result...)
	}

	return c.JSON(http.StatusOK, flightclass)
}

func (f *FlightController) FilterFlights(c echo.Context) error {
	airline := c.QueryParam("airline")
	aircraft := c.QueryParam("aircraft")
	departure := c.QueryParam("departure")
	capacity := c.QueryParam("capacity")

	url := fmt.Sprintf("%s/flights/filter?airline=%s&aircraft=%s&departure=%s&capacity=%s", utils.ENV("MOCK_URL"), airline, aircraft, departure, capacity)
	res, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get flights from mockapi",
		})
	}
	defer res.Body.Close()

	var flightclass []models.FlightClass
	err = json.NewDecoder(res.Body).Decode(&flightclass)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to decode flights from mockapi response",
		})
	}

	return c.JSON(http.StatusOK, flightclass)
}
