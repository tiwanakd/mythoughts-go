package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Username       string
	Email          string
	Name           string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email, name, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, name, hashed_password, created)
	VALUES ($1, $2, $3, $4, NOW())`

	_, err = m.DB.Exec(stmt, username, email, name, hashedPassword)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" && strings.Contains(pgErr.Message, "users_username_key") {
				return ErrDuplicateUsername
			}

			if pgErr.Code.Name() == "unique_violation" && strings.Contains(pgErr.Message, "users_email_key") {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, passoword string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentails
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(passoword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentails
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Get(id int) (User, error) {
	var user User
	stmt := "SELECT id, username, email, name, created FROM users WHERE id = $1"

	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}

		return User{}, err
	}

	return user, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) Update(id int, columnName, value string) (string, error) {

	query := fmt.Sprintf("UPDATE users SET %s = $1 WHERE id = $2 RETURNING %s", columnName, columnName)

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var newValue string
	err = stmt.QueryRow(value, id).Scan(&newValue)
	if err != nil {
		return "", err
	}

	return newValue, nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {

	stmt := "SELECT hashed_password FROM users WHERE id = $1"

	var hashedPassword []byte
	err := m.DB.QueryRow(stmt, id).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentails
		}
		return err
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	stmt = "UPDATE users SET hashed_password = $1 WHERE id = $2"
	_, err = m.DB.Exec(stmt, newHashedPassword, id)
	return err
}
