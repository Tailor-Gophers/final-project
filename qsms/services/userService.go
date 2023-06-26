package services

import (
	"qsms/models"
	"qsms/repository"
	"qsms/utils"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByUserName(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	DeleteUser(userId uint) error
	SaveToken(user *models.User, token string) error
	UserByToken(token string) (*models.User, error)
	LogOut(token string) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService() UserService {
	return &userService{
		userRepository: repository.NewGormUserRepository(),
	}
}

func (us *userService) CreateUser(user *models.User) error {
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return us.userRepository.CreateUser(user)
}

func (us *userService) GetUserByUserName(username string) (*models.User, error) {
	return us.userRepository.GetUserByUserName(username)
}

func (us *userService) GetUserByEmail(email string) (*models.User, error) {
	return us.userRepository.GetUserByEmail(email)
}

func (us *userService) SaveToken(user *models.User, token string) error {
	return us.userRepository.SaveToken(user, token)
}

func (us *userService) RefreshToken(user *models.User, token string) error {
	return us.userRepository.SaveToken(user, token)
}

func (us *userService) DeleteUser(userId uint) error {
	return us.userRepository.DeleteUser(userId)
}

func (us *userService) UserByToken(token string) (*models.User, error) {
	return us.userRepository.UserByToken(token)
}

func (us *userService) LogOut(token string) error {
	return us.userRepository.LogOut(token)
}