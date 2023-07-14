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
	RefreshToken(user *models.User, token string) error
	UserByToken(token string) (*models.User, error)
	LogOut(token string) error
	AddBalance(userId uint, amount int) error
	AddContact(user *models.User, contact *models.Contact) error
	DeleteContact(contactId uint) error
	GetContact(contactId uint) (*models.Contact, error)
	CreatePhoneBook(user *models.User, phonebook models.PhoneBook) error
	UpdatePhoneBook(phonebook *models.PhoneBook, number string) error
	GetPhoneBook(phonebookId uint) (*models.PhoneBook, error)
	DeletePhoneBook(phonebookId uint) error
	GetNumberByID(numberId uint) (*models.Number, error)
	GetUserByID(userId uint) (*models.User, error)
	GetAvailablePhoneNumbers() ([]models.Number, error)
	SetMainNumber(user *models.User, numberId uint) error
	CreateTemplate(template *models.Template) error
	DeleteTemplate(templateId uint) error
	GetTemplate(templateId uint) (*models.Template, error)
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
	return us.userRepository.UpdateBalance(userId, user.Balance+amount)
}

func (us *userService) AddContact(user *models.User, contact *models.Contact) error {
	return us.userRepository.AddContact(user, contact)
}

func (us *userService) CreateTemplate(template *models.Template) error {
	return us.userRepository.CreateTemplate(template)
}

func (us *userService) DeleteTemplate(templateId uint) error {
	return us.userRepository.DeleteTemplate(templateId)
}

func (us *userService) DeleteContact(contactId uint) error {
	return us.userRepository.DeleteContact(contactId)
}

func (us *userService) GetContact(contactId uint) (*models.Contact, error) {
	return us.userRepository.GetContact(contactId)
}

func (us *userService) CreatePhoneBook(user *models.User, phonebook models.PhoneBook) error {
	return us.userRepository.CreatePhoneBook(user, phonebook)
}

func (us *userService) UpdatePhoneBook(phonebook *models.PhoneBook, number string) error {
	return us.userRepository.UpdatePhoneBook(phonebook, number)
}

func (us *userService) GetPhoneBook(phonebookId uint) (*models.PhoneBook, error) {
	return us.userRepository.GetPhoneBook(phonebookId)
}

func (us *userService) GetNumberByID(numberId uint) (*models.Number, error) {
	return us.userRepository.GetNumberByID(numberId)
}

func (us *userService) DeletePhoneBook(phonebookId uint) error {
	return us.userRepository.DeletePhoneBook(phonebookId)
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

func (us *userService) GetTemplate(templateId uint) (*models.Template, error) {
	return us.userRepository.GetTemplate(templateId)
}
