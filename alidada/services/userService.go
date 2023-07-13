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
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService() UserService {
	return &userService{
		userRepository: repository.NewGormUserRepository(),
	}
}

func (s *userService) CreateUser(user *models.User) error {
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUserByUserName(username string) (*models.User, error) {
	return s.userRepository.GetUserByUserName(username)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *userService) CreatePassenger(passenger *models.Passenger) error {
	return s.userRepository.CreatePassenger(passenger)
}

func (s *userService) GetPassengers(user *models.User) ([]models.Passenger, error) {
	return s.userRepository.GetPassengers(user)
}

func (s *userService) GetMyTickets(user *models.User) ([]models.Reservation, error) {
	reservations, err := s.userRepository.GetReservationsByUserId(user.ID)

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

func (s *userService) GetMyTicketsPdf(user *models.User, id string) (string, error) {
	reservations, err := s.userRepository.GetReservationsByOrderId(id, user.ID)
	for i, _ := range reservations {
		reservations[i].FlightClass, err = GetFlightClassByID(int(reservations[i].FlightClassID))
		if err != nil {
			return "", err
		}
	}
	saveTo := fmt.Sprintf("pdf/ticketsOfOrder%s.pdf", id)

	GeneratePdf(reservations, saveTo)

	return saveTo, nil
}

// unit
func Sort(arr *[]models.CancellationCondition, start, end int) []models.CancellationCondition {
	if start < end {
		partitionIndex := partition(*arr, start, end)
		Sort(arr, start, partitionIndex-1)
		Sort(arr, partitionIndex+1, end)
	}
	return *arr
}

// unit
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
		fmt.Println(1)
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
		fmt.Println(2)

		return models.FlightClass{}, fmt.Errorf("redis error")

	}

	err = json.Unmarshal([]byte(val), &flightclass)
	if err != nil {
		return models.FlightClass{}, fmt.Errorf("Failed to decode flights from redis")
	}
	return flightclass, nil

}

func (s *userService) CancellTicket(user *models.User, id string) (string, error) {
	reservation, _ := s.userRepository.GetReservationById(id, user.ID)
	cancellationConditions, _ := s.userRepository.GetCancellationConditionsByFlightClassID(reservation.FlightClassID)
	sortedCancellationConditions := Sort(cancellationConditions, 0, len(*cancellationConditions)-1)
	flightclass, err := GetFlightClassByID(int(reservation.FlightClassID))
	if err != nil {
		errors.New("Failed to decode flights from mockapi")
	}
	penalty, err2 := PenaltyCalculation(reservation, &flightclass, sortedCancellationConditions)
	if err2 != nil {
		return "", err2
	}
	s.userRepository.AnnouncingcancellationToMockByFlightClassID(reservation.FlightClassID)
	s.userRepository.CancellReservationById(reservation.ID)
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
		col3 := fmt.Sprintf("https://Alidada.com/pass/reservation/%d", reservation.ID)
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
func (s *userService) SaveToken(user *models.User, token string) error {
	return s.userRepository.SaveToken(user, token)
}

func (s *userService) RefreshToken(user *models.User, token string) error {
	return s.userRepository.SaveToken(user, token)
}

func (s *userService) DeleteUser(username string) error {
	//todo if needed
	return nil
}

func (s *userService) UserByToken(token string) (*models.User, error) {
	return s.userRepository.UserByToken(token)
}

func (s *userService) LogOut(token string) error {
	return s.userRepository.LogOut(token)
}
