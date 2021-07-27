package model

import (
	"encoding/json"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type School struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Active   bool   `json:"active"`
	Finished bool   `json:"finished"`
}

func (s *School) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Title, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{4}\.[0-9]$`))),
		validation.Field(&s.Active, validation.By(notEqual(s.Finished))),
	)
}

func (s *School) ToString() string {
	str, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(str)
}
