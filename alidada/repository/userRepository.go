package repository

import (
	"alidada/db"
	"alidada/models"
	"alidada/utils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUserName(username string) (*models.User, error)
	GetUserByEmail(username string) (*models.User, error)
	GetUserByUserId(userId uint) (*models.User, error)
	DeleteUser(userId uint) error
	CreatePassenger(passenger *models.Passenger) error
	GetPassengers(user *models.User) ([]models.Passenger, error)
	SaveToken(user *models.User, token string) error
	UserByToken(token string) (*models.User, error)
	GetReservationsByUserId(id uint) ([]models.Reservation, error)
	GetReservationsByOrderId(id string, userId uint) ([]models.Reservation, error)
	GetReservationById(id string, userId uint) (*models.Reservation, error)
	GetCancellationConditionsByFlightClassID(id uint) (*[]models.CancellationCondition, error)
	AnnouncingcancellationToMockByFlightClassID(id uint) error
	CancellReservationById(id uint) error
	LogOut(token string) error
	PassReservation(id string) (string, error)
}

type userGormRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() UserRepository {
	return &userGormRepository{
		db: db.GetDbConnection(),
	}
}

func (ur *userGormRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userGormRepository) GetUserByUserName(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userGormRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userGormRepository) CreatePassenger(passenger *models.Passenger) error {
	var findpassenger models.Passenger
	err := ur.db.Where("passport_number = ?", passenger.PassportNumber).First(&findpassenger).Error
	if err != nil {
		return ur.db.Create(passenger).Error
	}
	return fmt.Errorf("Passenger alredy exist")
}

func (ur *userGormRepository) GetPassengers(user *models.User) ([]models.Passenger, error) {
	var passengers []models.Passenger
	err := ur.db.Model(user).Association("Passengers").Find(&passengers)
	if err != nil {
		return nil, err
	}
	return passengers, nil
}

func (ur *userGormRepository) GetReservationById(id string, userId uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := ur.db.
		Joins("JOIN passengers ON passengers.id = reservations.passenger_id AND passengers.user_id = ?", userId).
		Where("confirmed =?", true).
		Where("reservations.is_cancelled != 1").
		Where("reservations.id = ?", id).
		First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (ur *userGormRepository) GetCancellationConditionsByFlightClassID(id uint) (*[]models.CancellationCondition, error) {
	var cancellationConditions []models.CancellationCondition
	err := ur.db.
		Joins("JOIN flight_class_cancellations ON flight_class_cancellations.cancellation_condition_id = cancellation_conditions.id").
		Where("flight_class_cancellations.flight_class_id = ?", id).
		Find(&cancellationConditions).Error
	if err != nil {
		return nil, err
	}
	return &cancellationConditions, nil
}

func (ur *userGormRepository) AnnouncingcancellationToMockByFlightClassID(id uint) error {
	url := fmt.Sprintf("%s/flights/%d/return", utils.ENV("MOCK_URL"), id)
	res, err := http.Post(url, "", nil)
	if err != nil {
		return fmt.Errorf("Failed to decode flights from mockapi")
	}
	defer res.Body.Close()
	return nil
}
func (ur *userGormRepository) CancellReservationById(id uint) error {
	return ur.db.Model(&models.Reservation{}).Where("id = ?", id).Update("is_cancelled", true).Error
}

func (ur *userGormRepository) GetReservationsByUserId(id uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := ur.db.Joins("JOIN passengers ON passengers.id = reservations.passenger_id AND passengers.user_id = ?", id).
		Select("reservations.*").Where("confirmed =?", true).
		Preload("Passenger").Find(&reservations).Error

	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (ur *userGormRepository) GetReservationsByOrderId(id string, userId uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := ur.db.Joins("JOIN passengers ON passengers.id = reservations.passenger_id AND passengers.user_id = ?", userId).
		Select("reservations.ID", "price", "passenger_id", "flight_class_id").
		Where("order_id = ?", id).
		Where("confirmed =?", true).
		Where("is_cancelled != 1").
		Preload("Passenger").
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	if len(reservations) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return reservations, nil
}

func (ur *userGormRepository) PassReservation(id string) (string, error) {
	var reservations []models.Reservation
	var reservation models.Reservation

	err := ur.db.First(&reservation, id).Error
	if err != nil {
		return "", err
	}
	if reservation.IsCancelled == true || reservation.Confirmed == false {
		return "", fmt.Errorf("this reservation Is Cancelled of not confirmed")
	}
	err = ur.db.
		Where("flight_class_id = ?", reservation.FlightClassID).
		Where("confirmed =?", true).
		Where("is_cancelled != 1").Where("id <= ?", id).
		Preload("Passenger").
		Find(&reservations).Error
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d_%d", int(reservation.FlightClassID), int(len(reservations)+1)), nil
}

func (ur *userGormRepository) DeleteUser(userId uint) error {
	return ur.db.Delete(&models.User{}, userId).Error
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

func (ur *userGormRepository) SaveToken(user *models.User, token string) error {

	hashed, err := utils.HashToken(token)
	if err != nil {
		return err
	}
	AccessToken := models.Token{UserId: user.ID, Token: hashed, ExpiresAt: time.Now().Add(time.Hour * 24)}

	result := ur.db.Create(&AccessToken)

	return result.Error
}

func (ur *userGormRepository) UserByToken(token string) (*models.User, error) {
	var AccessToken models.Token
	var User models.User

	hashed, err := utils.HashToken(token)
	if err != nil {
		return nil, err
	}
	err = ur.db.Where("token = ?", hashed).Where("expires_at > ?", time.Now()).First(&AccessToken).Error
	if err != nil {
		return nil, err
	}
	err = ur.db.Preload("Passengers").Where("id = ?", AccessToken.UserId).First(&User).Error
	if err != nil {
		return nil, err
	}
	return &User, nil
}

func (ur *userGormRepository) LogOut(token string) error {
	var AccessToken models.Token

	hashed, err := utils.HashToken(token)
	if err != nil {
		return err
	}
	err = ur.db.Where("token = ?", hashed).Where("expires_at > ?", time.Now()).First(&AccessToken).Error

	ur.db.Where("token = ?", hashed).Where("expires_at > ?", time.Now()).Delete(&AccessToken)
	if err != nil {
		return err
	}
	return nil
}
