package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Leson struct {
	ID    int64
	Title string
}

func (l *Leson) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Title, validation.Required),
	)
}
