package services

import (
	"errors"

	"lab6server/internal/domen/models"
)

type UsersStorage interface {
	Add(user *models.User) error
	Delete(login string) error
	Get(login string) (*models.User, error)
	Update(user *models.User) error
}
type AuthService struct {
	userStorage UsersStorage
}

func NewAuthService(userStorage UsersStorage) *AuthService {
	return &AuthService{
		userStorage: userStorage,
	}
}

func (a *AuthService) Authenticate(login, password string) error {

	user, err := a.userStorage.Get(login)
	if err != nil {
		err = a.register(&models.User{Login: login, Password: password})
		if err != nil {
			return err
		}
		return nil
	}

	if user.Password != password {
		return errors.New("incorrect password")
	}

	return nil
}

func (a *AuthService) register(user *models.User) error {

	err := a.userStorage.Add(user)
	if err != nil {
		return errors.New("failed to register user")
	}

	return nil
}
