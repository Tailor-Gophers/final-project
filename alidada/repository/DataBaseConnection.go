package repository

import (
	"final-project/alidada/models"
	"final-project/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getDbConnection() *gorm.DB {

	dbURI := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s?charset=utf8&parseTime=True&loc=Local", utils.ENV("DB_USERNAME"), utils.ENV("DB_PASSWORD"), utils.ENV("DB_PORT"), utils.ENV("DB_DATABASE"))
	// Connect to the database

	db, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.AccessToken{}, &models.User{})
	if err != nil {
		panic(err)
	}
	return db
}
