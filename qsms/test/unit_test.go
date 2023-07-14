package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"qsms/models"
	"qsms/repository"
	"qsms/services"
	"testing"
	"time"
)

func TestGenerateTextFromTemplate(t *testing.T) {
	user := &models.User{
		UserName: "JohnDoe",
	}
	template := "Hello, {{user_name}}! Today is {{date}} at {{time}}."

	expected := "Hello, JohnDoe! Today is " + time.Now().Format("2006-01-02") + " at " + time.Now().Format("15:04:5") + "."

	result := services.GenerateTextFromTemplate(user, template)
	assert.Equal(t, expected, result)
}

func TestSearchMessages(t *testing.T) {
	words := []string{"Hello", "world"}
	messages := []models.Message{
		{ID: 1, Message: "Hello there!"},
		{ID: 2, Message: "Howdy, world!"},
		{ID: 3, Message: "Hi, everyone!"},
	}

	adminService := services.NewAdminService(&mockAdminRepository{messages: messages})

	result, err := adminService.SearchMessages(words)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, messages[0], result[0])
	assert.Equal(t, messages[1], result[1])
}

func (mar *mockAdminRepository) GetAllMessages() ([]models.Message, error) {
	return mar.messages, nil
}

type mockAdminRepository struct {
	mock.Mock
	repository.AdminRepository
	messages []models.Message
}
