package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"qsms/models"
	"qsms/repository"
	"qsms/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

type MessageService interface {
	SendSimpleMessage(user *models.User, receiver string, text string) error
	SendTemplateMessage(user *models.User, receiver string, template string) error
	SendPeriodicSimpleMessage(user *models.User, receiver string, text string, interval string) error
	SendPeriodicTemplateMessage(user *models.User, receiver string, template string, interval string) error
	RegisterMessagingSchedules() error
}

type messageService struct {
	MessageRepository repository.MessageRepository
	UserRepository    repository.UserRepository
}

type Words struct {
	Words []string `json:"words"`
}

var regex *regexp.Regexp
var config *utils.Config

func NewMessageService(messageRepository repository.MessageRepository,
	userRepository repository.UserRepository) MessageService {

	config = utils.LoadConfig()

	pattern := "\\b(" + strings.Join(config.BadWords, "|") + ")\\b"
	regex = regexp.MustCompile(pattern)

	return &messageService{
		MessageRepository: messageRepository,
		UserRepository:    userRepository,
	}
}

func (ms *messageService) SendSimpleMessage(user *models.User, receiver string, text string) error {

	if user.Balance < config.SimpleMessageFee {
		return errors.New("insufficient balance")
	}
	err := ms.UserRepository.UpdateBalance(user.ID, user.Balance-config.SimpleMessageFee)
	if err != nil {
		return err
	}

	if CheckBadWords(text) {
		return errors.New("can't send message: contains bad words")
	}
	if strings.TrimSpace(text) == "" {
		return errors.New("message can't be empty")
	}

	log.Printf("Message from %s(%d) -> %s :: %s \n", user.UserName, user.MainNumberID, receiver, text)
	err = ms.sendSms(receiver, text)
	if err != nil {
		return errors.New("failed to send sms: " + err.Error())
	}

	message := &models.Message{
		SenderID:       user.ID,
		ReceiverNumber: receiver,
		Message:        text,
	}
	err = ms.MessageRepository.SaveMessage(message)
	if err != nil {
		return errors.New("failed to save message: " + err.Error())
	}
	return nil
}

func (ms *messageService) SendTemplateMessage(user *models.User, receiver string, template string) error {

	totalCost := config.SimpleMessageFee + config.TemplateFee
	if user.Balance < totalCost {
		return errors.New("insufficient balance")
	}
	err := ms.UserRepository.UpdateBalance(user.ID, user.Balance-totalCost)
	if err != nil {
		return err
	}

	text := GenerateTextFromTemplate(user, template)

	if CheckBadWords(text) {
		return errors.New("can't send message: contains bad words")
	}

	log.Printf("Message from %s(%d) -> %s :: %s \n", user.UserName, user.MainNumberID, receiver, text)
	err = ms.sendSms(receiver, text)
	if err != nil {
		return errors.New("failed to send sms: " + err.Error())
	}

	message := &models.Message{
		SenderID:       user.ID,
		ReceiverNumber: receiver,
		Message:        text,
	}
	err = ms.MessageRepository.SaveMessage(message)
	if err != nil {
		return errors.New("failed to save message: " + err.Error())
	}
	return nil
}

func (ms *messageService) SendPeriodicSimpleMessage(user *models.User, receiver string, text string, interval string) error {

	totalCost := config.SimpleMessageFee + config.PeriodicMessageFee //having money for at least one message sending

	if user.Balance < totalCost {
		return errors.New("insufficient balance")
	}
	err := ms.UserRepository.UpdateBalance(user.ID, user.Balance-config.PeriodicMessageFee)
	if err != nil {
		return errors.New("failed to update balance: " + err.Error())
	}

	schedule := &models.MessageSchedule{
		UserID:   user.ID,
		Receiver: receiver,
		Text:     text,
		Interval: interval,
	}
	err = ms.MessageRepository.SaveScheduler(schedule)
	if err != nil {
		return errors.New("failed to save scheduler: " + err.Error())
	}

	s := ScheduleWithParser(interval)
	_, err = s.Do(ms.SendSimpleMessage, user, receiver, text)
	if err != nil {
		return errors.New("failed to schedule task: " + err.Error())
	}
	s.StartAsync()
	return nil
}

func (ms *messageService) SendPeriodicTemplateMessage(user *models.User, receiver string, template string, interval string) error {
	totalCost := config.PeriodicMessageFee + config.PeriodicMessageFee + config.TemplateFee //having money for at least one message sending

	if user.Balance < totalCost {
		return errors.New("insufficient balance")
	}
	err := ms.UserRepository.UpdateBalance(user.ID, user.Balance-(config.PeriodicMessageFee+config.TemplateFee))
	if err != nil {
		return err
	}

	schedule := &models.MessageSchedule{
		UserID:   user.ID,
		Receiver: receiver,
		Template: template,
		Interval: interval,
	}
	err = ms.MessageRepository.SaveScheduler(schedule)
	if err != nil {
		return errors.New("failed to save scheduler: " + err.Error())
	}

	s := ScheduleWithParser(interval)
	_, err = s.Do(ms.SendTemplateMessage, user, receiver, template)
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}

func (ms *messageService) RegisterMessagingSchedules() error {
	schedules, err := ms.MessageRepository.GetAllSchedules()
	if err != nil {
		return err
	}

	for _, schedule := range schedules {
		s := ScheduleWithParser(schedule.Interval)
		user, err := ms.UserRepository.GetUserById(schedule.UserID)
		if err != nil {
			return errors.New("user not found: " + err.Error())
		}
		if schedule.Text != "" { //is simple
			_, err = s.Do(ms.SendSimpleMessage, user, schedule.Receiver, schedule.Text)
			if err != nil {
				return errors.New("failed to schedule task: " + err.Error())
			}
		} else if schedule.Template != "" {
			_, err = s.Do(ms.SendTemplateMessage, user, schedule.Receiver, schedule.Template)
			if err != nil {
				return errors.New("failed to schedule task: " + err.Error())
			}
		}
		s.StartAsync()
	}
	return nil
}

func GenerateTextFromTemplate(user *models.User, template string) string {
	//can add more
	replacements := map[string]string{
		"{{user_name}}": user.UserName,
		"{{date}}":      time.Now().Format("2006-01-02"),
		"{{time}}":      time.Now().Format("15:04:5"),
		"{{date_time}}": time.Now().String(),
	}
	for v, s := range replacements {
		template = strings.Replace(template, v, s, -1)
	}
	return template
}

func ScheduleWithParser(interval string) *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)
	switch string(interval[len(interval)-1]) {
	case "s":
		return s.Every(interval)
	case "m":
		return s.Every(interval)
	case "h":
		return s.Every(interval)
	case "d":
		numb, _ := strconv.Atoi(interval[:len(interval)-1])
		return s.Every(numb).Days()
	case "M":
		numb, _ := strconv.Atoi(interval[:len(interval)-1])
		return s.Every(24 * numb).Days()
	}
	return s
}

func CheckBadWords(text string) bool {
	if regex.MatchString(text) {
		return true
	}
	return false
}

type MockResponse struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`
}

func (ms *messageService) sendSms(receiver string, text string) error {
	postBody, _ := json.Marshal(map[string]string{
		"receiver": receiver,
		"message":  text,
	})
	resp, err := http.Post(utils.ENV("MOCK_URL")+"/send", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return errors.New("failed to send request: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response MockResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("error from mock: " + response.ErrorMessage)
	}
	return nil
}
