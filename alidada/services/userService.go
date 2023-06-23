package services

import (
	"alidada/models"
	"alidada/repository"
	"alidada/utils"
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
	GetMyTicketsPdf(user *models.User, id string) ([]models.Reservation, error)
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
	return s.userRepository.GetMyTickets(user)
}
func (s *userService) GetMyTicketsPdf(user *models.User, id string) ([]models.Reservation, error) {
	return s.userRepository.GetMyTicketsPdf(user, id)
}

func (s *userService) CancellTicket(user *models.User, id string) (string, error) {
	return s.userRepository.CancellTicket(user, id)
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
