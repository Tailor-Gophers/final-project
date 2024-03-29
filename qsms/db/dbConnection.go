package db

import (
	"fmt"
	"qsms/models"
	"qsms/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDbConnection() *gorm.DB {

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", utils.ENV("DB_USERNAME"), utils.ENV("DB_PASSWORD"), utils.ENV("DB_URL"), utils.ENV("DB_PORT"), utils.ENV("DB_DATABASE"))

	db, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.User{}, &models.PhoneBook{}, &models.RawNumber{}, &models.Contact{}, &models.Number{}, &models.Token{}, &models.Transaction{}, &models.Rent{}, &models.Template{}, &models.Message{}, &models.MessageSchedule{})
	if err != nil {
		panic(err)
	}

	return db
}
