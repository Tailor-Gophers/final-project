package repository

import (
	"gorm.io/gorm"
	"qsms/db"
	"qsms/models"
)

type PaymentRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransaction(authority string) (models.Transaction, error)
	ConfirmTransaction(transaction *models.Transaction) error
}

type paymentGormRepository struct {
	db *gorm.DB
}

func NewGormPaymentRepository() PaymentRepository {
	return &paymentGormRepository{
		db: db.GetDbConnection(),
	}
}

func (pr *paymentGormRepository) CreateTransaction(transaction *models.Transaction) error {
	return pr.db.Create(transaction).Error
}

func (pr *paymentGormRepository) GetTransaction(authority string) (models.Transaction, error) {
	var transaction models.Transaction
	err := pr.db.Where("authority = ?", authority).First(&transaction).Error
	return transaction, err
}

func (pr *paymentGormRepository) ConfirmTransaction(transaction *models.Transaction) error {
	return pr.db.Save(transaction).Error
}
