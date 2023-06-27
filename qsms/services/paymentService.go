package services

import (
	"qsms/models"
	"qsms/repository"
)

type PaymentService interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransaction(authority string) (models.Transaction, error)
	ConfirmTransaction(transaction *models.Transaction, refId int) error
}

type paymentService struct {
	paymentRepository repository.PaymentRepository
}

func NewPaymentService(repository repository.PaymentRepository) PaymentService {
	return &paymentService{
		paymentRepository: repository,
	}
}

func (ps *paymentService) CreateTransaction(transaction *models.Transaction) error {
	return ps.paymentRepository.CreateTransaction(transaction)
}

func (ps *paymentService) GetTransaction(authority string) (models.Transaction, error) {
	return ps.paymentRepository.GetTransaction(authority)
}

func (ps *paymentService) ConfirmTransaction(transaction *models.Transaction, refId int) error {
	transaction.Confirmed = true
	transaction.RefID = refId
	return ps.paymentRepository.ConfirmTransaction(transaction)
}
