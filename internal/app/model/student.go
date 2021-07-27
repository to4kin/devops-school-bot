package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Student struct {
	ID      int64    `json:"id"`
	Account *Account `json:"account"`
	School  *School  `json:"school"`
	Active  bool     `json:"active"`
}

func (s *Student) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Account, validation.Required),
		validation.Field(&s.School, validation.Required),
	)
}

func (s *Student) ToString() string {
	str, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(str)
}
