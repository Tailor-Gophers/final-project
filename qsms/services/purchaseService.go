package services

import (
	"errors"
	"fmt"
	"log"
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

const RentingRate = 0.1

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

	ps.ScheduleRentingTask(24*30*time.Hour, user.ID, number.ID)

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
		err = ps.DropRent(user, rent.ID)
		if err != nil {
			log.Println("Error while dropping rent in scheduler!")
		}
		return errors.New(fmt.Sprintf("Low balance: you need at least %d to rent number %d", rent.Price, numberId))
	}

	err = ps.purchaseRepository.SetUserBalance(user, user.Balance-rent.Price)
	if err != nil {
		return err
	}
	err = ps.purchaseRepository.UpdateRentDate(rent.ID, time.Now())
	if err != nil {
		fmt.Println(2)
		return err
	}

	ps.ScheduleRentingTask(24*30*time.Hour, user.ID, numberId)

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

	number, err := ps.purchaseRepository.GetNumberByID(rent.NumberID)
	if err != nil {
		return err
	}

	if user.MainNumberID == number.ID {
		err = ps.purchaseRepository.UpdateUserMainNumber(user.ID)
		if err != nil {
			return err
		}
	}

	err = ps.purchaseRepository.RestoreNumber(number.ID)
	if err != nil {
		return err
	}
	err = ps.purchaseRepository.DropRent(rentId)
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

	for _, rent := range rents {
		delay := rent.LastPaid.Add(24 * 30 * time.Hour).Sub(time.Now())
		ps.ScheduleRentingTask(delay, rent.UserID, rent.NumberID)
	}
	return nil
}

func (ps *purchaseService) ScheduleRentingTask(delay time.Duration, userId uint, numberId uint) {
	go func() {
		time.Sleep(delay)
		err := ps.UpdateRent(userId, numberId)
		if err != nil {
			log.Printf("Failed to update rent for user %d and number %d: "+err.Error(), userId, numberId)
		}
	}()
}
