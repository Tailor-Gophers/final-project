package repository

import (
	"final-project/alidada/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getDbConnection() *gorm.DB {
	//user:password@/dbname?charset=utf8&parseTime=True&loc=Local")

	dbURI := "admin:13771377Ab?@tcp(localhost:3306)/quera?charset=utf8&parseTime=True&loc=Local"
	// Connect to the database

	db, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Set up connection pool and other configuration options

	// Enable logging in development mode
	// Migrate the User model to the database (if necessary)
	db.AutoMigrate(&models.AccessToken{}, &models.User{})
	// db.Model(&models.Carpet{}).AddForeignKey("collection_id", "collections(id)", "RESTRICT", "RESTRICT")
	// db.Model(&models.CarpetColor{}).AddForeignKey("carpet_id", "carpets(id)", "RESTRICT", "RESTRICT")
	// db.Model(&models.CarpetMedia{}).AddForeignKey("carpet_color_id", "carpet_colors(id)", "RESTRICT", "RESTRICT")
	// db.Model(&models.CarpetMedia{}).AddForeignKey("carpet_id", "carpets(id)", "RESTRICT", "RESTRICT")

	// Use the db instance to interact with the database in your application
	return db
}
