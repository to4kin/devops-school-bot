package model

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type School struct {
	ID         int64
	Title      string
	InProgress bool
	Finished   bool
}

func (s *School) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Title, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{4}\.[0-9]$`))),
		validation.Field(&s.InProgress, validation.By(notEqual(s.Finished))),
	)
}
