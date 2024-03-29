package db

import (
	"alidada/models"
	"alidada/utils"
	"fmt"

	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDbConnection() *gorm.DB {

	dbURI := fmt.Sprintf("%s:%s@tcp(mysql1:%s)/%s?charset=utf8&parseTime=True&loc=Local", utils.ENV("DB_USERNAME"), utils.ENV("DB_PASSWORD"), utils.ENV("DB_PORT"), utils.ENV("DB_DATABASE"))
	// Connect to the database

	db, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.CancellationCondition{}, &models.User{}, &models.Passenger{}, &models.Token{}, &models.Order{}, &models.Reservation{}, &models.AuthorityPair{}, &models.FlightClassCancellation{})
	if err != nil {
		panic(err)
	}

	return db
}

func GetRedisConnection() *redis.Client {
	redisURI := fmt.Sprintf("redis:%s", utils.ENV("REDIS_PORT"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: utils.ENV("REDIS_PASSWORD"),
		DB:       1,
	})

	return rdb

}
