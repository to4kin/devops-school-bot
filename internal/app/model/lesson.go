package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Lesson ...
type Lesson struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// GetID ...
func (l *Lesson) GetID() int64 {
	return l.ID
}

// GetButtonTitle ...
func (l *Lesson) GetButtonTitle() string {
	return l.Title
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

// LogrusFields ...
func (l *Lesson) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":    l.ID,
		"title": l.Title,
	}
}
