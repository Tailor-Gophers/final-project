package controllers

import (
	"mockapi/services"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type FlightController struct {
	FlightService services.FlightService
}

func (f *FlightController) GetFlightByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	result, err := f.FlightService.GetFlight(int64(id))
	if err != nil {
		return c.String(http.StatusNotFound, "Flight Not Found!!")

	}
	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetFlightByDate(c echo.Context) error {

	origin := c.Param("origin")
	destination := c.Param("destination")
	dateStr := c.Param("date")

	day, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid date format")
	}

	result, err := f.FlightService.GetFlightByDate(origin, destination, day)
	if err != nil {
		return c.String(http.StatusNotFound, "Flight Not Found!!")
	}

	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetPlanesList(c echo.Context) error {
	result, err := f.FlightService.GetPlanesList()
	if err != nil {
		return c.String(http.StatusNotFound, "No Planes Found!!")

	}
	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetCitiesList(c echo.Context) error {
	result, err := f.FlightService.GetCitiesList()
	if err != nil {
		return c.String(http.StatusNotFound, "No City Found!!")

	}
	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetDaysList(c echo.Context) error {
	result, err := f.FlightService.GetDaysList()
	if err != nil {
		return c.String(http.StatusNotFound, "No Day Found!!")

	}
	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) ReserveFlightCapacity(c echo.Context) error { // Reduce Capacity
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	result, err := f.FlightService.ReserveFlightCapacity(int64(id))
	if err != nil {
		if err.Error() == "flight capacity reached" {
			return c.String(http.StatusNotFound, "Flight capacity reached!")
		}
		return c.String(http.StatusNotFound, "No Flight Found!")
	}
	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) ReturnFlightCapacity(c echo.Context) error { // Increase Capacity
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	result, err := f.FlightService.ReturnFlightCapacity(int64(id))
	if err != nil {
		if err.Error() == "flight capacity is already empty" {
			return c.String(http.StatusNotFound, "Flight capacity is already empty!")
		}
		return c.String(http.StatusNotFound, "No Flight Found!")
	}
	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetFlightByFilter(c echo.Context) error {
	airline := c.Param("airline")
	aircraft := c.Param("aircraft")
	departureStr := c.Param("departure")

	departure, err := time.Parse("2006-01-02", departureStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid date format")
	}

	result, err := f.FlightService.GetFlightByFilter(airline, aircraft, departure)
	if err != nil {
		return c.String(http.StatusNotFound, "Flight Not Found!!")
	}

	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetFlightBySort(c echo.Context) error {
	order := c.Param("order")
	result, err := f.FlightService.GetFlightBySort(order)
	if err != nil {
		return c.String(http.StatusNotFound, "Flight Not Found!!")
	}

	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetFlightPrice(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	result, err := f.FlightService.GetFlightPrice(int64(id))
	if err != nil {
		return c.String(http.StatusNotFound, "Flight Not Found!!")
	}

	return c.JSON(http.StatusOK, result)
}
