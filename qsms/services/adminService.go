package services

import (
	"qsms/models"
	"qsms/repository"
	"regexp"
)

type AdminService interface {
	AddNumber(number *models.Number) error
	SuspendUser(userId uint) error
	UnSuspendUser(userId uint) error
	CountUserMessages(userId uint) (int, error)
	SearchMessages(words []string) ([]models.Message, error)
	GetMessageByID(messageId uint) (models.Message, error)
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

func (as *adminService) CountUserMessages(userId uint) (int, error) {
	return as.AdminRepository.CountUserMessages(userId)
}

func (as *adminService) SearchMessages(words []string) ([]models.Message, error) {
	messages, err := as.AdminRepository.GetAllMessages()
	if err != nil {
		return nil, err
	}

	pattern := "\\b(" + regexp.QuoteMeta(words[0])
	for _, word := range words[1:] {
		pattern += "|" + regexp.QuoteMeta(word)
	}
	pattern += ")\\b"
	regex := regexp.MustCompile(pattern)

	result := make([]models.Message, 0)

	for _, message := range messages {
		if regex.MatchString(message.Message) {
			result = append(result, AdminViewMessage(message))
		}
	}
	return messages, nil
}

func (as *adminService) GetMessageByID(messageId uint) (models.Message, error) {
	message, err := as.AdminRepository.GetMessageByID(messageId)
	if err != nil {
		return models.Message{}, err
	}
	return AdminViewMessage(message), nil
}

func AdminViewMessage(message models.Message) models.Message {
	pattern := `\b\d{5,7}\b`
	regex := regexp.MustCompile(pattern)
	maskedInput := regex.ReplaceAllString(message.Message, "******")

	message.Message = maskedInput

	return message
}
