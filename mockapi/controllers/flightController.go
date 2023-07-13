package controllers

import (
	"mockapi/models"
	"mockapi/services"
	"net/http"
	"strconv"
	"time"

	echo "github.com/labstack/echo/v4"
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
	origin, destination, err := f.FlightService.GetCitiesList()
	if err != nil {
		return c.String(http.StatusNotFound, "No City Found!!")
	}

	cities := map[string]interface{}{
		"Origin":      origin,
		"Destination": destination,
	}

	return c.JSON(http.StatusOK, cities)
}

func (f *FlightController) GetDaysList(c echo.Context) error {
	result, err := f.FlightService.GetDaysList()
	if err != nil {
		return c.String(http.StatusNotFound, "No Day Found!!")
	}

	responseData := map[string][]models.Flight{}

	for _, flight := range result {
		day := flight.StartTime.Format("2006-01-02")
		responseData[day] = append(responseData[day], flight)
	}

	return c.JSON(http.StatusOK, responseData)
}

func (f *FlightController) ReserveFlightCapacity(c echo.Context) error { // Reduce Capacity
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	result, err := f.FlightService.ReserveFlightCapacity(int64(id), int(count))
	if err != nil {
		if err.Error() == "flight capacity reached" {
			return c.String(http.StatusNotFound, "Flight capacity reached!")
		}
		return c.String(http.StatusNotFound, err.Error())
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
	airline := c.QueryParam("airline")
	aircraft := c.QueryParam("aircraft")
	departureStr := c.QueryParam("departure")

	capacityInt, err := strconv.Atoi(c.QueryParam("capacity"))
	if err != nil && c.QueryParam("capacity") != "" {
		return c.String(http.StatusBadRequest, "Invalid capacity")
	}
	capacity := uint(capacityInt)

	departure, err := time.Parse("2006-01-02", departureStr)
	if err != nil && c.QueryParam("departure") != "" {
		return c.String(http.StatusBadRequest, "Invalid date format")
	}

	result, err := f.FlightService.GetFlightByFilter(airline, aircraft, departure, capacity)
	if err != nil {
		return c.String(http.StatusNotFound, "Flight Not Found!!")
	}

	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetFlightBySort(c echo.Context) error {
	orderBy := c.Param("orderBy")
	order := c.QueryParam("order")
	result, err := f.FlightService.GetFlightBySort(orderBy, order)
	if err != nil {
		return c.String(http.StatusNotFound, "Flight Not Found!!")
	}

	return c.JSON(http.StatusOK, result)
}

func (f *FlightController) GetFlightClassByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	result, err := f.FlightService.GetFlightClassByID(int64(id))
	if err != nil {
		return c.String(http.StatusNotFound, "Flight CLass Not Found!!")
	}

	return c.JSON(http.StatusOK, result)
}
