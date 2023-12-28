package models

import (
	"database/sql"
	"time"
)

type Users struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Created      time.Time
}

type UsersModel struct {
	DB *sql.DB
}

func (user *UsersModel) Insert(name, email, password string) (int, error) {
	return 0, nil
}

func (user *UsersModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (user *UsersModel) Exists(name, email string) (bool, error) {
	return false, nil
}
