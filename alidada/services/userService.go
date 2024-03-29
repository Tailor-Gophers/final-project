package services

import (
	"alidada/db"
	"alidada/models"
	"alidada/repository"
	"alidada/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"

	redis "github.com/redis/go-redis/v9"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByUserName(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	DeleteUser(username string) error
	CreatePassenger(passenger *models.Passenger) error
	GetPassengers(user *models.User) ([]models.Passenger, error)
	SaveToken(user *models.User, token string) error
	UserByToken(token string) (*models.User, error)
	LogOut(token string) error
	GetMyTickets(user *models.User) ([]models.Reservation, error)
	CancellTicket(user *models.User, id string) (string, error)
	GetMyTicketsPdf(user *models.User, id string) (string, error)
	PassReservation(id string) (string, error)
}

type UserServicet struct {
	UserRepository repository.UserRepository
}

func NewUserService() UserService {
	return &UserServicet{
		UserRepository: repository.NewGormUserRepository(),
	}
}

func (s *UserServicet) CreateUser(user *models.User) error {
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.UserRepository.CreateUser(user)
}

func (s *UserServicet) GetUserByUserName(username string) (*models.User, error) {
	return s.UserRepository.GetUserByUserName(username)
}

func (s *UserServicet) GetUserByEmail(email string) (*models.User, error) {
	return s.UserRepository.GetUserByEmail(email)
}

func (s *UserServicet) CreatePassenger(passenger *models.Passenger) error {
	return s.UserRepository.CreatePassenger(passenger)
}

func (s *UserServicet) GetPassengers(user *models.User) ([]models.Passenger, error) {
	return s.UserRepository.GetPassengers(user)
}

func (s *UserServicet) GetMyTickets(user *models.User) ([]models.Reservation, error) {
	reservations, err := s.UserRepository.GetReservationsByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	for i, _ := range reservations {
		reservations[i].FlightClass, err = GetFlightClassByID(int(reservations[i].FlightClassID))
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (s *UserServicet) GetMyTicketsPdf(user *models.User, id string) (string, error) {
	reservations, err := s.UserRepository.GetReservationsByOrderId(id, user.ID)
	if err != nil {
		return "", err
	}
	for i, _ := range reservations {
		reservations[i].FlightClass, err = GetFlightClassByID(int(reservations[i].FlightClassID))
		if err != nil {
			return "", err
		}
	}
	saveTo := fmt.Sprintf("pdf/ticketsOfOrder%s.pdf", id)

	err = GeneratePdf(reservations, saveTo)
	if err != nil {
		return "", err
	}

	return saveTo, nil
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

// unit
func PenaltyCalculation(reservation *models.Reservation, flightclass *models.FlightClass, sortedCancellationConditions []models.CancellationCondition) (int, error) {
	for _, condition := range sortedCancellationConditions {
		var t time.Duration
		t = time.Duration(condition.TimeMinutes)
		if flightclass.Flight.StartTime.Unix() < time.Now().Add(-1*time.Minute*t).Unix() {
			return condition.Penalty * int(reservation.Price) / 100, nil
		}
	}
	return 100, errors.New("None of the cancellation conditions are available for you")
}

func GetFlightClassByID(id int) (models.FlightClass, error) {
	ctx := context.Background()
	rdb := db.GetRedisConnection()
	key := fmt.Sprintf("flightClass_%d", id)
	val, err := rdb.Get(ctx, key).Result()
	var flightclass models.FlightClass

	if err == redis.Nil {
		//key does not exist
		url := fmt.Sprintf("%s/flight_class/%d", utils.ENV("MOCK_URL"), id)

		res, err := http.Get(url)
		if err != nil {
			return models.FlightClass{}, fmt.Errorf("Failed to decode flights from mockapi")
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&flightclass)
		if err != nil {
			return models.FlightClass{}, fmt.Errorf("Failed to decode flights from mockapi")
		}

		flightclassMarshal, _ := json.Marshal(flightclass)

		err2 := rdb.Set(ctx, key, flightclassMarshal, 100000000000).Err()
		if err2 != nil {
			return models.FlightClass{}, fmt.Errorf("cant Saving to redis")
		}
		return flightclass, nil
	}
	if err != nil {

		return models.FlightClass{}, fmt.Errorf("redis error")

	}

	err = json.Unmarshal([]byte(val), &flightclass)
	if err != nil {
		return models.FlightClass{}, fmt.Errorf("Failed to decode flights from redis")
	}
	return flightclass, nil

}

func (s *UserServicet) CancellTicket(user *models.User, id string) (string, error) {
	reservation, err := s.UserRepository.GetReservationById(id, user.ID)
	if err != nil {
		return "", err
	}
	cancellationConditions, err := s.UserRepository.GetCancellationConditionsByFlightClassID(reservation.FlightClassID)
	if err != nil {
		return "", err
	}
	sortedCancellationConditions := Sort(cancellationConditions, 0, len(*cancellationConditions)-1)
	flightclass, err := GetFlightClassByID(int(reservation.FlightClassID))
	if err != nil {
		errors.New("Failed to decode flights from mockapi")
	}
	penalty, err := PenaltyCalculation(reservation, &flightclass, sortedCancellationConditions)
	if err != nil {
		return "", err
	}
	err = s.UserRepository.AnnouncingcancellationToMockByFlightClassID(reservation.FlightClassID)
	if err != nil {
		return "", err
	}
	err = s.UserRepository.CancellReservationById(reservation.ID)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("your penalty is: %d", penalty)
	return result, nil
}

// unit
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
		col3 := fmt.Sprintf("%s/api/pass/reservation/%d", utils.ENV("URL"), reservation.ID)
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
		return err2
	}
	return nil
}
func (s *UserServicet) SaveToken(user *models.User, token string) error {
	return s.UserRepository.SaveToken(user, token)
}

func (s *UserServicet) RefreshToken(user *models.User, token string) error {
	return s.UserRepository.SaveToken(user, token)
}

func (s *UserServicet) DeleteUser(username string) error {
	//todo if needed
	return nil
}

func (s *UserServicet) UserByToken(token string) (*models.User, error) {
	return s.UserRepository.UserByToken(token)
}

func (s *UserServicet) LogOut(token string) error {
	return s.UserRepository.LogOut(token)
}

func (s *UserServicet) PassReservation(id string) (string, error) {
	return s.UserRepository.PassReservation(id)
}
