package usecase

import (
	"time"

	"lab6server/internal/domen/models"
)

type AuthService interface {
	Authenticate(login string, password string) error
}

type MessageStorage interface {
	Add(message *models.Message)
	GetMessages() []*models.Message
}

type ChatUsecase struct {
	authService    AuthService
	messageStorage MessageStorage
}

func NewChatUsecase(authService AuthService, messageStorage MessageStorage) *ChatUsecase {
	return &ChatUsecase{
		authService:    authService,
		messageStorage: messageStorage,
	}
}

func (uc *ChatUsecase) GetChatHistory() ([]*models.Message, error) {
	history := uc.messageStorage.GetMessages()
	return history, nil
}

func (uc *ChatUsecase) AddMessage(login, password, text string) (*models.Message, error) {

	err := uc.authService.Authenticate(login, password)
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		Username: login,
		Text:     text,
		Date:     time.Now(),
	}

	uc.messageStorage.Add(message)
	return message, nil
}
