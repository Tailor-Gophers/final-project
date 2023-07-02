package services

import (
	"fmt"
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
	AddBalance(userId uint, amount int) error
	AddContact(contact *models.Contact) error
	DeleteContact(user *models.User, contactId uint) error
	GetContact(contactId uint) (*models.Contact, error)
	UpdateContact(user *models.User, contact *models.Contact) error
	GetUserByID(userId uint) (*models.User, error)
	GetAvailablePhoneNumbers() ([]models.Number, error)
	SetMainNumber(user *models.User, numberId uint) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		userRepository: repository,
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

func (us *userService) AddBalance(userId uint, amount int) error {
	user, err := us.userRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	fmt.Println(user.Balance + amount)
	return us.userRepository.UpdateBalance(userId, user.Balance+amount)
}

func (us *userService) AddContact(contact *models.Contact) error {
	return us.userRepository.AddContact(contact)
}

func (us *userService) DeleteContact(user *models.User, contactId uint) error {
	return us.userRepository.DeleteContact(user, contactId)
}

func (us *userService) UpdateContact(user *models.User, contact *models.Contact) error {
	return us.userRepository.UpdateContact(user, contact)
}

func (us *userService) GetContact(contactId uint) (*models.Contact, error) {
	return us.userRepository.GetContact(contactId)
}

func (us *userService) GetUserByID(userId uint) (*models.User, error) {
	return us.userRepository.GetUserById(userId)
}

func (us *userService) GetAvailablePhoneNumbers() ([]models.Number, error) {
	return us.userRepository.GetAvailablePhoneNumbers()
}

func (us *userService) SetMainNumber(user *models.User, numberId uint) error {
	return us.userRepository.SetMainNumber(user, numberId)
}
