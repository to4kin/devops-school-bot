package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Homework struct {
	ID      int64
	Student *Student
	Lesson  *Lesson
	Accept  bool
}

func (h *Homework) Validate() error {
	return validation.ValidateStruct(
		h,
		validation.Field(&h.Student, validation.Required),
		validation.Field(&h.Lesson, validation.Required),
	)
}
