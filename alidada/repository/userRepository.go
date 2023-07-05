package repository

import (
	"alidada/db"
	"alidada/models"
	"alidada/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
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
	LogOut(token string) error
	GetMyTickets(user *models.User) ([]models.Reservation, error)
	GetFlightClassByID(id int) (models.FlightClass, error)
	CancellTicket(user *models.User, id string) (string, error)
	GetMyTicketsPdf(user *models.User, id string) (string, error)
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
	return ur.db.Create(passenger).Error
}

func (ur *userGormRepository) GetPassengers(user *models.User) ([]models.Passenger, error) {
	var passengers []models.Passenger
	err := ur.db.Model(user).Association("Passengers").Find(&passengers)
	if err != nil {
		return nil, err
	}
	return passengers, nil
}

func Sort(arr *[]models.CancellationCondition, start, end int) []models.CancellationCondition {
	if start < end {
		partitionIndex := partition(*arr, start, end)
		Sort(arr, start, partitionIndex-1)
		Sort(arr, partitionIndex+1, end)
	}
	return *arr
}

func partition(arr []models.CancellationCondition, start, end int) int {
	pivot := arr[end].Penalty
	pIndex := start
	for i := start; i < end; i++ {
		if arr[i].Penalty <= pivot {
			//  swap
			arr[i], arr[pIndex] = arr[pIndex], arr[i]
			pIndex++
		}
	}
	arr[pIndex], arr[end] = arr[end], arr[pIndex]
	return pIndex
}

func (ur *userGormRepository) PenaltyCalculation(reservation *models.Reservation) (int, error) {
	var cancellationConditions []models.CancellationCondition
	ur.db.
		Joins("JOIN flight_class_cancellations ON flight_class_cancellations.cancellation_condition_id = cancellation_conditions.id").
		Where("flight_class_cancellations.flight_class_id = ?", reservation.FlightClassID).
		Find(&cancellationConditions)

	sortedCancellationConditions := Sort(&cancellationConditions, 0, len(cancellationConditions)-1)
	for _, condition := range sortedCancellationConditions {
		var t time.Duration
		t = time.Duration(condition.TimeMinutes)
		if reservation.CreatedAt.Unix() < time.Now().Add(-1*time.Minute*t).Unix() {
			return condition.Penalty * int(reservation.Price) / 100, nil
		}
	}
	return 100, errors.New("None of the cancellation conditions are available for you")
}

func (ur *userGormRepository) CancellTicket(user *models.User, id string) (string, error) {
	var reservation models.Reservation
	err := ur.db.
		Joins("JOIN passengers ON passengers.id = reservations.passenger_id AND passengers.user_id = ?", user.ID).
		Where("reservations.is_cancelled is null").
		Where("reservations.id = ?", id).
		First(&reservation).Error
	if err != nil {
		return "", err
	}

	penalty, err2 := ur.PenaltyCalculation(&reservation)
	if err2 != nil {
		return "", err2
	}

	url := fmt.Sprintf("http://localhost:3001/flights/%d/return", reservation.FlightClassID)
	res, err := http.Post(url, "", nil)
	if err != nil {
		return "", fmt.Errorf("Failed to decode flights from mockapi")
	}
	defer res.Body.Close()
	ur.db.Model(&models.Reservation{}).Where("id = ?", reservation.ID).Update("is_cancelled", true)

	result := fmt.Sprintf("your penalty is: %d", penalty)

	return result, nil
}
func (ur *userGormRepository) GetMyTickets(user *models.User) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := ur.db.Joins("JOIN passengers ON passengers.id = reservations.passenger_id AND passengers.user_id = ?", user.ID).
		Select("*").
		Preload("Passenger").Find(&reservations).Error

	for i, _ := range reservations {
		reservations[i].FlightClass, err = ur.GetFlightClassByID(int(reservations[i].FlightClassID))
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (ur *userGormRepository) GetFlightClassByID(id int) (models.FlightClass, error) {
	url := fmt.Sprintf("http://localhost:3001/flight_class/%d", id)
	res, err := http.Get(url)
	if err != nil {
		return models.FlightClass{}, fmt.Errorf("Failed to decode flights from mockapi")
	}
	defer res.Body.Close()

	var flightclass models.FlightClass
	err = json.NewDecoder(res.Body).Decode(&flightclass)
	if err != nil {
		return models.FlightClass{}, fmt.Errorf("Failed to decode flights from mockapi")
	}

	return flightclass, nil
}

func GeneratePdf(reservations []models.Reservation, saveTo string) error {
	m := pdf.NewMaroto(consts.Portrait, consts.Letter)
	//m.SetBorder(true)
	m.AddUTF8Font("Shabnam", consts.Normal, "Shabnam.ttf")
	for i, reservation := range reservations {

		m.Row(40, func() {
			m.Col(4, func() {
				_ = m.FileImage("static/airplane.png", props.Rect{
					Center:  true,
					Percent: 80,
				})
			})

			m.Col(4, func() {
				m.Text(" Ali Dada Airlines | Tailor Gopher, Inc. ", props.Text{
					Top:         12,
					Size:        20,
					Family:      "Shabnam",
					Extrapolate: true,
				})

				m.Text("Automatic ticket issuing system", props.Text{
					Size: 12,
					Top:  22,
				})
			})
			m.ColSpace(4)
		})

		m.Line(10)
		col1 := fmt.Sprintf("%d- Name: %s %s | Date of birth: %s | National code: %s | Passport : %s ", i+1, reservation.Passenger.FirstName, reservation.Passenger.LastName, reservation.Passenger.DateOfBirth, reservation.Passenger.Nationality, reservation.Passenger.PassportNumber)
		col2 := fmt.Sprintf("https://Alidada.com/passenger/%d", reservation.PassengerID)
		col3 := fmt.Sprintf("https://Alidada.com/reservation/%d", reservation.ID)
		col4 := fmt.Sprintf("CODE: %d | DATE: %s | Origin: %s | Destination: %s ", reservation.FlightClass.ID, reservation.FlightClass.Flight.StartTime.Format("2006-01-02 15:04:05"), reservation.FlightClass.Flight.Origin, reservation.FlightClass.Flight.Destination)

		m.Row(40, func() {
			m.Col(4, func() {
				m.Text(col1, props.Text{
					Size:   15,
					Top:    12,
					Family: "Shabnam",
				})
			})
			m.ColSpace(4)
			m.Col(4, func() {
				m.QrCode(col2, props.Rect{
					Center:  true,
					Percent: 75,
				})
			})
		})

		m.Line(10)

		m.Row(100, func() {
			m.Col(12, func() {
				_ = m.Barcode(col3, props.Barcode{
					Center:  true,
					Percent: 70,
				})
				m.Text("AliDada . com", props.Text{
					Size:  20,
					Align: consts.Center,
					Top:   65,
				})
			})
		})

		m.SetBorder(true)

		m.Row(40, func() {
			m.Col(6, func() {
				m.Text(col4, props.Text{
					Size: 15,
					Top:  14,
				})
			})
			m.Col(6, func() {
				m.Text(reservation.FlightClass.Title, props.Text{
					Top:   1,
					Size:  50,
					Align: consts.Center,
				})
			})
		})

		m.SetBorder(false)

	}
	err2 := m.OutputFileAndClose(saveTo)
	if err2 != nil {
		return fmt.Errorf("pdf cant build")
	}
	return nil
}

func (ur *userGormRepository) GetMyTicketsPdf(user *models.User, id string) (string, error) {
	var reservations []models.Reservation
	err := ur.db.Joins("JOIN passengers ON passengers.id = reservations.passenger_id AND passengers.user_id = ?", user.ID).
		Select("reservations.ID", "price", "passenger_id", "flight_class_id").
		Where("order_id = ?", id).
		Preload("Passenger").
		Find(&reservations).Error
	if err != nil {
		return "", err
	}

	if len(reservations) == 0 {
		return "", fmt.Errorf("not found")
	}

	for i, _ := range reservations {
		reservations[i].FlightClass, err = ur.GetFlightClassByID(int(reservations[i].FlightClassID))
		if err != nil {
			return "", err
		}
	}
	saveTo := fmt.Sprintf("pdf/ticketsOfOrder%s.pdf", id)

	GeneratePdf(reservations, saveTo)

	return saveTo, nil
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
