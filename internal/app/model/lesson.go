package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Lesson struct {
	ID    int64
	Title string
}

func (l *Lesson) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Title, validation.Required),
	)
}
