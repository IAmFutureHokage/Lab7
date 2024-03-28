package messagestorage

import (
	"lab6server/internal/domen/models"
	"sort"
	"sync"
	"time"
)

type MessageStorage struct {
	mu       sync.RWMutex
	messages []*models.Message
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		messages: make([]*models.Message, 0),
	}
}

func (s *MessageStorage) Add(message *models.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, message)
}

func (s *MessageStorage) GetMessages() []*models.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messagesCopy := make([]*models.Message, len(s.messages))
	copy(messagesCopy, s.messages)
	return messagesCopy
}

func (s *MessageStorage) RemoveOldMessages(olderThan time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := sort.Search(len(s.messages), func(i int) bool {
		return !s.messages[i].Date.Before(olderThan)
	})

	if index >= len(s.messages) {
		s.messages = []*models.Message{}
		return
	}
	s.messages = s.messages[index:]
}
