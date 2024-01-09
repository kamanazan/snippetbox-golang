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

func (user *UsersModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO users(name, email, password_hash, created)
	VALUES($1, $2, $3, localtimestamp);
	`

	_, err = user.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		/*
			23505	unique_violation
		*/
		var pgErr *pq.Error

		if errors.As(err, &pgErr) {
			if strings.Contains(pgErr.Message, "users_email_key") {
				fmt.Printf("duplicate error %s", pgErr)
				return ErrDuplicateEmail
			} else {
				fmt.Printf("unknown error %s", pgErr)
				return pgErr
			}
		}
	}

	return nil
}

func (user *UsersModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (user *UsersModel) Exists(name, email string) (bool, error) {
	return false, nil
}
