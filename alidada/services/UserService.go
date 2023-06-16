package services

import (
	"final-project/alidada/models"
	"final-project/alidada/repository"
	"final-project/utils"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByUserName(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	DeleteUser(username string) error
	SaveToken(user *models.User, token string) error
	UserByToken(token string) (*models.User, error)
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

func (s *userService) GetUserByEmail(username string) (*models.User, error) {
	return s.userRepository.GetUserByEmail(username)
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
