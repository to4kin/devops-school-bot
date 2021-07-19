package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Student ...
type Student struct {
	ID         int
	TelegramID int
	FirstName  string
	LastName   string
	Username   string
}

func (s *Student) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.TelegramID, validation.Required),
		validation.Field(&s.FirstName, validation.Required),
	)
}
