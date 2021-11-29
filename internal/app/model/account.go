package model

import (
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Account object reprsents a student
type Account struct {
	// ID returns an Account.ID
	//
	// NOTE: filled in automatically after INSERT to the store
	ID int64 `json:"id,string"`

	// Created returns time.Time
	//
	// NOTE: should be set before INSRT to the store
	// *field is required
	Created time.Time `json:"created"`

	// TelegramID returns user telegram id
	//
	// NOTE: should be equal to real user telegram id for future search
	// *field is required
	TelegramID int64 `json:"telegram_id,string"`

	// FirstName returns user first name
	//
	// NOTE: doesn't participate in search, so could be different then in telegram
	// *field is required
	FirstName string `json:"first_name"`

	// LastName returns user last name
	//
	// NOTE: doesn't participate in search, so could be different then in telegram
	LastName string `json:"last_name"`

	// Username returns telegram username
	//
	// NOTE: should be equal to real telegram username for @ in reports
	Username string `json:"username"`

	// Superuser returns true/false
	//
	// NOTE: false by default if not specified
	Superuser bool `json:"superuser,string" `
}

// GetID Account.ID
func (a *Account) GetID() int64 {
	return a.ID
}

// GetButtonTitle returns "@" + Account.Username
func (a *Account) GetButtonTitle() string {
	return "@" + a.Username
}

// GetFullName returns a full name for account
//
// NOTE:
// - returns FirstName + LastName by default
// - returns FirstName only if LastName is empty
// - returns LastName only if FirstName is empty
func (a *Account) GetFullName() string {
	if a.FirstName == "" {
		return a.LastName
	}

	if a.LastName == "" {
		return a.FirstName
	}

	return fmt.Sprintf("%v %v", a.FirstName, a.LastName)
}

// GetMention returns mention of a user
//
// NOTE: mention is returned in HTML style
// <a href="tg://user?id=Account.TelegramID">Account.GetFullName()</a>
func (a *Account) GetMention() string {
	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%v</a>", a.TelegramID, a.GetFullName())
}

// Validate func is needed to validate Account object fields before INSERT
//
// NOTE:
// - Created is required
// - TelegramID is requiredq
// - FirstName is required
func (a *Account) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Created, validation.Required),
		validation.Field(&a.TelegramID, validation.Required),
		validation.Field(&a.FirstName, validation.Required),
	)
}

// ToString converts Account object to json string
func (a *Account) ToString() string {
	str, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields returns logrus.Fields for logrus logger
//
// NOTE:
// available fields are "id", "telegram_id", "first_name", "last_name",
// "username", "superuser"
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
