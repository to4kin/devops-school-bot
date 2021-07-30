package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Student ...
type Student struct {
	ID      int64    `json:"id"`
	Account *Account `json:"account"`
	School  *School  `json:"school"`
	Active  bool     `json:"active"`
}

// Validate ...
func (s *Student) Validate() error {
	return validation.ValidateStruct(
		s,
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
