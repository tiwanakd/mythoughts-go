package mocks

import (
	"time"

	"github.com/tiwanakd/mythoughts-go/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(username, email, name, password string) error {
	if username == "dupeUsername" {
		return models.ErrDuplicateUsername
	}

	if email == "dupe@email.com" {
		return models.ErrDuplicateEmail
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "test@test.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentails
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (models.User, error) {
	switch id {
	case 1:
		return models.User{
			ID:             1,
			Username:       "test",
			Email:          "test@test.com",
			Name:           "test user",
			HashedPassword: []byte("pa$$word"),
			Created:        time.Now(),
		}, nil
	default:
		return models.User{}, models.ErrInvalidCredentails
	}
}

func (m *UserModel) Update(id int, columnName, value string) error {
	return nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	return nil
}

func (m *UserModel) Delete(id int) error {
	switch id {
	case 1:
		return nil
	default:
		return models.ErrInvalidCredentails
	}
}
