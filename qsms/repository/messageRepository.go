package repository

import (
	"gorm.io/gorm"
	"qsms/db"
	"qsms/models"
)

type MessageRepository interface {
	SaveMessage(message *models.Message) error
	SaveScheduler(schedule *models.MessageSchedule) error
	GetAllSchedules() ([]models.MessageSchedule, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewGormMessageRepository() MessageRepository {
	return &messageRepository{
		db: db.GetDbConnection(),
	}
}

func (mr *messageRepository) SaveMessage(message *models.Message) error {
	return mr.db.Create(message).Error
}

func (mr *messageRepository) SaveScheduler(schedule *models.MessageSchedule) error {
	return mr.db.Create(schedule).Error
}

func (mr *messageRepository) GetAllSchedules() ([]models.MessageSchedule, error) {
	var schedules []models.MessageSchedule
	err := mr.db.Find(&schedules).Error
	return schedules, err
}
