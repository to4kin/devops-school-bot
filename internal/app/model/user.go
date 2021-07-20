package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID         int64
	TelegramID int64
	FirstName  string
	LastName   string
	Username   string
	IsAdmin    bool
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.TelegramID, validation.Required),
		validation.Field(&u.FirstName, validation.Required),
	)
}
