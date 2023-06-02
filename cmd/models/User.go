package models

import (
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/leomarzochi/facebooklike/cmd/crypt"
)

const (
	USER_STATUS_CREATING = "creating"
	USER_STATUS_EDITING  = "updating"
)

type User struct {
	ID       uint64    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	CreateAt time.Time `json:"createdAt,omitempty"`
}

// Will validate and format received user
func (u *User) Prepare(status string) error {
	if err := u.validate(status); err != nil {
		return err
	}

	err := u.format(status)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) validate(status string) error {
	if u.Name == "" {
		return errors.New("field name is required")
	}

	if u.Username == "" {
		return errors.New("field username is required")
	}

	if u.Email == "" {
		return errors.New("field email is required")
	}

	err := checkmail.ValidateFormat(u.Email)
	if err != nil {
		return errors.New("invalid e-mail format")
	}

	if status == "creating" && u.Password == "" {
		return errors.New("field password is required")
	}

	return nil
}

func (u *User) format(status string) error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.Username = strings.TrimSpace(u.Username)

	if status == USER_STATUS_CREATING {
		hashedPassword, err := crypt.HashPassword(u.Password)
		if err != nil {
			return err
		}

		u.Password = string(hashedPassword)
	}

	return nil
}
