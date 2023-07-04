package services

import (
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"qsms/models"
	"qsms/repository"
	"time"
)

type PurchaseService interface {
	BuyNumber(user *models.User, numberId uint) error
	RegisterRentingTasks() error
	PlaceRent(user *models.User, numberId uint) error
	UpdateRent(userId, numberId uint) error
	DropRent(user *models.User, rentId uint) error
}

type purchaseService struct {
	purchaseRepository repository.PurchaseRepository
}

func NewPurchaseService(repository repository.PurchaseRepository) PurchaseService {
	return &purchaseService{
		purchaseRepository: repository,
	}
}

var RentingRate = 0.1

func (ps *purchaseService) BuyNumber(user *models.User, numberId uint) error {
	number, err := ps.purchaseRepository.GetNumberByID(numberId)
	if err != nil {
		return err
	}

	if number.Active == true {
		return errors.New(fmt.Sprintf("Number with id %d is already owned!", numberId))
	}

	if user.Balance < number.Price {
		return errors.New("insufficient balance")
	}

	err = ps.purchaseRepository.SetUserBalance(user, user.Balance-number.Price)
	if err != nil {
		return err
	}

	err = ps.purchaseRepository.UpdateNumber(user, numberId)
	if err != nil {
		return err
	}

	return nil
}

func (ps *purchaseService) RegisterRentingTasks() error {
	rents, err := ps.purchaseRepository.GetAllRents()
	if err != nil {
		return err
	}

	s := gocron.NewScheduler(time.UTC)
	for _, rent := range rents {
		_, err = s.Every(30).Days().
			StartAt(rent.LastPaid.Add(30*24*time.Hour)).Do(ps.UpdateRent, rent.UserID, rent.NumberID)
		if err != nil {
			return err
		}
	}
	s.StartAsync()
	return nil
}

func (ps *purchaseService) PlaceRent(user *models.User, numberId uint) error {
	number, err := ps.purchaseRepository.GetNumberByID(numberId)
	if err != nil {
		return err
	}

	if number.Active == true {
		return errors.New(fmt.Sprintf("Number with id %d is not available!", numberId))
	}

	rentPrice := int(float64(number.Price) * RentingRate)

	if user.Balance < rentPrice {
		return errors.New(
			fmt.Sprintf("Low balance: you need at least %d to rent number %d", rentPrice, numberId))
	}

	err = ps.purchaseRepository.SetUserBalance(user, user.Balance-rentPrice)
	if err != nil {
		return err
	}

	err = ps.purchaseRepository.UpdateNumber(user, number.ID)
	if err != nil {
		return err
	}

	rent := &models.Rent{
		UserID:   user.ID,
		NumberID: number.ID,
		Price:    rentPrice,
		LastPaid: time.Now(),
	}
	err = ps.purchaseRepository.PlaceRent(rent)
	if err != nil {
		return err
	}

	s := gocron.NewScheduler(time.UTC)
	scheduleTime := time.Now().Add(30 * 24 * time.Hour)
	_, err = s.Every(30).Days().StartAt(scheduleTime).Do(ps.UpdateRent, user.ID, number.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.StartAsync()
	return nil
}

func (ps *purchaseService) UpdateRent(userId, numberId uint) error {
	rent, err := ps.purchaseRepository.GetRentByIDs(userId, numberId)
	if err != nil {
		return err
	}
	user, err := ps.purchaseRepository.GetUserById(userId)
	if err != nil {
		return err
	}

	if user.Disable == true {
		return errors.New("user account is currently disabled")
	}
	if user.Balance < rent.Price {
		return errors.New(fmt.Sprintf("Low balance: you need at least %d to rent number %d", rent.Price, numberId))
	}

	err = ps.purchaseRepository.SetUserBalance(user, user.Balance-rent.Price)
	if err != nil {
		return err
	}
	err = ps.purchaseRepository.UpdateRentDate(rent.ID, time.Now())
	if err != nil {
		return err
	}

	s := gocron.NewScheduler(time.UTC)
	scheduleTime := time.Now().Add(30 * 24 * time.Hour)
	_, err = s.Every(30).Days().StartAt(scheduleTime).Do(ps.UpdateRent, user.ID, numberId)
	if err != nil {
		return err
	}
	s.StartAsync()

	return nil
}

func (ps *purchaseService) DropRent(user *models.User, rentId uint) error {
	rent, err := ps.purchaseRepository.GetRentByID(rentId)
	if err != nil {
		return err
	}

	if user.ID != rent.UserID {
		return errors.New(fmt.Sprintf("User has no renting record with id %d!", rentId))
	}
	fmt.Println("2")
	number, err := ps.purchaseRepository.GetNumberByID(rent.NumberID)
	if err != nil {
		return err
	}
	fmt.Println("3")
	err = ps.purchaseRepository.RestoreNumber(number.ID)
	if err != nil {
		return err
	}
	fmt.Println("4")
	err = ps.purchaseRepository.DropRent(rentId)
	if err != nil {
		return err
	}
	fmt.Println("5")
	return nil
}
