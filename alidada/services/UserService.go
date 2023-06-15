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
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
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

func (s *userService) DeleteUser(username string) error {
	//todo if needed
	return nil
}
