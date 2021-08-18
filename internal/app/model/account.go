package model

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Account ...
type Account struct {
	ID         int64     `json:"id,string"`
	Created    time.Time `json:"created"`
	TelegramID int64     `json:"telegram_id,string"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	Superuser  bool      `json:"superuser,string"`
}

// GetID ...
func (a *Account) GetID() int64 {
	return a.ID
}

// GetButtonTitle ...
func (a *Account) GetButtonTitle() string {
	return "@" + a.Username
}

// Validate ...
func (a *Account) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Created, validation.Required),
		validation.Field(&a.TelegramID, validation.Required),
		validation.Field(&a.FirstName, validation.Required),
	)
}

// ToString ...
func (a *Account) ToString() string {
	str, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields ...
func (a *Account) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":          a.ID,
		"telegram_id": a.TelegramID,
		"first_name":  a.FirstName,
		"last_name":   a.LastName,
		"username":    a.Username,
		"superuser":   a.Superuser,
	}
}
