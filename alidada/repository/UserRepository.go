package repository

import (
	"errors"
	"final-project/alidada/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUserId(userId uint) (*models.User, error)
	DeleteUser(user *models.User) error
}

type userGormRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() UserRepository {
	return &userGormRepository{
		db: getDbConnection(),
	}
}

func (ur *userGormRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userGormRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
func (ur *userGormRepository) DeleteUser(user *models.User) error {
	return ur.db.Delete(user).Error
}

func (ur *userGormRepository) GetUserByUserId(userId uint) (*models.User, error) {
	var user models.User
	result := ur.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

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
