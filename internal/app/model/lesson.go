package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Lesson ...
type Lesson struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// Validate ...
func (l *Lesson) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Title, validation.Required),
	)
}

// ToString ...
func (l *Lesson) ToString() string {
	str, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}
