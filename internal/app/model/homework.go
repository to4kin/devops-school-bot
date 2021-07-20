package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Homework struct {
	ID    int64
	Title string
}

func (h *Homework) Validate() error {
	return validation.ValidateStruct(
		h,
		validation.Field(&h.Title, validation.Required),
	)
}
