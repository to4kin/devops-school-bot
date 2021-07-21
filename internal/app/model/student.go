package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Student struct {
	ID      int64
	Account *Account
	School  *School
	Active  bool
}

func (s *Student) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Account, validation.Required),
		validation.Field(&s.School, validation.Required),
	)
}
