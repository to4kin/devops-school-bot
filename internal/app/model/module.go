package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Module ...
type Module struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// GetID ...
func (l *Module) GetID() int64 {
	return l.ID
}

// GetButtonTitle ...
func (l *Module) GetButtonTitle() string {
	return l.Title
}

// Validate ...
func (l *Module) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Title, validation.Required),
	)
}

// ToString ...
func (l *Module) ToString() string {
	str, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields ...
func (l *Module) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":    l.ID,
		"title": l.Title,
	}
}
