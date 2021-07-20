package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Account struct {
	ID         int64
	TelegramID int64
	FirstName  string
	LastName   string
	Username   string
	IsAdmin    bool
}

func (a *Account) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.TelegramID, validation.Required),
		validation.Field(&a.FirstName, validation.Required),
	)
}
