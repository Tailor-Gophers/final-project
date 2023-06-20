package controllers

import (
	"alidada/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FlightController struct{}

func NewFlightController() *FlightController {
	return &FlightController{}
}

func (f *FlightController) SearchFlights(c echo.Context) error {
	origin := c.QueryParam("origin")
	destination := c.QueryParam("destination")
	dateStr := c.QueryParam("date")

	url := fmt.Sprintf("http://localhost:3001/flights/%s/%s/%s", origin, destination, dateStr)
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
