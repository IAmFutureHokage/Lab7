package usersstorage

import (
	"errors"
	"lab6server/internal/domen/models"
	"sync"
)

type UserStorage struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[string]*models.User),
	}
}

func (s *UserStorage) Add(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Login]; exists {
		return errors.New("user already exists")
	}

	s.users[user.Login] = user
	return nil
}

func (s *UserStorage) Get(login string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[login]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserStorage) Delete(login string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[login]; !exists {
		return errors.New("user not found")
	}

	delete(s.users, login)
	return nil
}

func (s *UserStorage) Update(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Login]; !exists {
		return errors.New("user not found")
	}

	s.users[user.Login] = user
	return nil
}
