package db

import (
	"fmt"
	"mockapi/models"
	"mockapi/utils"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDbConnection() *gorm.DB {

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", utils.ENV("DB_USERNAME"), utils.ENV("DB_PASSWORD"), utils.ENV("DB_URL"), utils.ENV("DB_PORT"), utils.ENV("DB_DATABASE"))
	// Connect to the database

	db, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Flight{}, &models.FlightClass{})
	if err != nil {
		panic(err)
	}

	var count int64
	db.Model(&models.Flight{}).Count(&count)
	if count == 0 {
		formatTime := "2006-01-02 15:04:05"
		startTime1, _ := time.Parse(formatTime, "2020-02-09 14:49:41")
		endTime1, _ := time.Parse(formatTime, "2020-02-09 18:49:41")
		startTime2, _ := time.Parse(formatTime, "2020-02-09 19:49:41")
		endTime2, _ := time.Parse(formatTime, "2020-02-09 20:49:41")
		startTime3, _ := time.Parse(formatTime, "2020-02-08 15:49:41")
		endTime3, _ := time.Parse(formatTime, "2020-02-09 18:49:41")
		startTime4, _ := time.Parse(formatTime, "2020-02-13 09:30:41")
		endTime4, _ := time.Parse(formatTime, "2020-02-13 10:00:41")

		flights := []models.Flight{
			{Model: models.Model{ID: 1}, Origin: "Tehran", Destination: "Tabriz", StartTime: startTime1, EndTime: endTime1, Airline: "homa", Aircraft: "Boeing737"},
			{Model: models.Model{ID: 3}, Origin: "Tehran", Destination: "Mashad", StartTime: startTime2, EndTime: endTime2, Airline: "homa", Aircraft: "Boeing737"},
			{Model: models.Model{ID: 2}, Origin: "Shiraz", Destination: "Tehran", StartTime: startTime3, EndTime: endTime3, Airline: "homa", Aircraft: "Boeing737"},
			{Model: models.Model{ID: 6}, Origin: "Tehran", Destination: "Alborz", StartTime: startTime4, EndTime: endTime4, Airline: "mahan", Aircraft: "Boeing737"},
		}

		for _, flight := range flights {
			db.Create(&flight)
		}

		var reserve1, reserve2, reserve3, reserve4 uint = 15, 15, 1, 49
		flightClasses := []models.FlightClass{
			{Model: models.Model{ID: 4}, Title: "Class-A", Price: 1900, Capacity: 50, Reserve: &reserve1, FlightId: 2},
			{Model: models.Model{ID: 1}, Title: "Class-A", Price: 1900, Capacity: 50, Reserve: &reserve2, FlightId: 6},
			{Model: models.Model{ID: 2}, Title: "Class-B", Price: 1700, Capacity: 50, Reserve: &reserve3, FlightId: 6},
			{Model: models.Model{ID: 3}, Title: "Class-C", Price: 1300, Capacity: 50, Reserve: &reserve4, FlightId: 6},
		}

		for _, flightClass := range flightClasses {
			db.Create(&flightClass)
		}
	}

	return db
}
