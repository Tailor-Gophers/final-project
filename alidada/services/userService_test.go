package services

import (
	"alidada/models"
	"alidada/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockUserRepository struct {
	mock.Mock
	repository.UserRepository
}

func (m *mockUserRepository) CreatePassenger(passenger *models.Passenger) error {
	args := m.Called(passenger)
	return args.Error(0)
}

func (m *mockUserRepository) GetPassengers(user *models.User) ([]models.Passenger, error) {
	args := m.Called(user)
	return args.Get(0).([]models.Passenger), args.Error(1)
}

func TestCreatePassenger(t *testing.T) {
	userRepo := &mockUserRepository{}
	service := userService{userRepository: userRepo}

	mockPassenger := &models.Passenger{
		Model: models.Model{
			ID: 2,
		},
		UserID:         2,
		FirstName:      "Behzad",
		LastName:       "Ebram",
		Gender:         false,
		DateOfBirth:    "000-1",
		Nationality:    "Iranian",
		PassportNumber: "123456712",
	}

	userRepo.On("CreatePassenger", mockPassenger).Return(nil)

	err := service.CreatePassenger(mockPassenger)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestGetPassengers(t *testing.T) {
	userRepo := &mockUserRepository{}
	service := userService{userRepository: userRepo}

	mockUser := &models.User{
		Model: models.Model{
			ID: 2,
		},
		Email:       "amir@me.com",
		UserName:    "amir",
		Password:    "amir1234",
		FirstName:   "Amir",
		LastName:    "Val",
		PhoneNumber: "1234567890",
		IsAdmin:     false,
		Passengers:  []models.Passenger{},
		Tokens:      []models.Token{},
	}

	userRepo.On("GetPassengers", mockUser).Return([]models.Passenger{}, nil)

	_, err := service.GetPassengers(mockUser)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}
