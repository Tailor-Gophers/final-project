package services

import (
	"qsms/models"
	"qsms/repository"
)

type AdminService interface {
	AddNumber(number *models.Number) error
	SuspendUser(userId uint) error
	UnSuspendUser(userId uint) error
}

type adminService struct {
	AdminRepository repository.AdminRepository
}

func NewAdminService(repository repository.AdminRepository) AdminService {
	return &adminService{
		AdminRepository: repository,
	}
}

func (as *adminService) AddNumber(number *models.Number) error {
	return as.AdminRepository.AddNumber(number)
}

func (as *adminService) SuspendUser(userId uint) error {
	return as.AdminRepository.SuspendUser(userId)
}

func (as *adminService) UnSuspendUser(userId uint) error {
	return as.AdminRepository.UnSuspendUser(userId)
}
