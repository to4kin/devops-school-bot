package model

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Student ...
type Student struct {
	ID      int64     `json:"id"`
	Created time.Time `json:"created"`
	Account *Account  `json:"account"`
	School  *School   `json:"school"`
	Active  bool      `json:"active"`
}

// Validate ...
func (s *Student) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Created, validation.Required),
		validation.Field(&s.Account, validation.Required),
		validation.Field(&s.School, validation.Required),
	)
}

// ToString ...
func (s *Student) ToString() string {
	str, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields ...
func (s *Student) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":      s.ID,
		"account": s.Account.LogrusFields(),
		"school":  s.School.LogrusFields(),
		"active":  s.Active,
	}
}
