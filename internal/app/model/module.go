package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Module object represents a module
type Module struct {
	// ID returns a Module.ID
	//
	// NOTE: filled in automatically after INSERT to the store
	ID int64 `json:"id"`

	// Title returns a module title
	//
	// NOTE:
	// *field is required
	Title string `json:"title"`
}

// GetID returns Module.ID
func (l *Module) GetID() int64 {
	return l.ID
}

// GetButtonTitle returns Module.Title
func (l *Module) GetButtonTitle() string {
	return l.Title
}

// Validate func is needed to validate Module object before INSERT
//
// NOTE:
// - Title is required
func (l *Module) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Title, validation.Required),
	)
}

// ToString converts Module object to json string
func (l *Module) ToString() string {
	str, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields returns logrus.Fields for logrus logger
//
// NOTE:
// available fields are "id", "title"
func (l *Module) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":    l.ID,
		"title": l.Title,
	}
}
