package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Account struct {
	ID         int64  `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
	Superuser  bool   `json:"superuser"`
}

func (a *Account) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.TelegramID, validation.Required),
		validation.Field(&a.FirstName, validation.Required),
	)
}

func (a *Account) ToString() string {
	str, err := json.Marshal(a)
	if err != nil {
		return ""
	}

	return string(str)
}
