package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Lesson object represents a lesson
type Lesson struct {
	// ID returns a Lesson.ID
	//
	// NOTE: filled in automatically after INSERT to the store
	ID int64 `json:"id"`

	// Title returns a lesson title
	//
	// NOTE:
	// *field is required
	Title string `json:"title"`

	// Module object
	//
	// NOTE:
	// *object is required
	Module *Module `json:"module"`
}

// GetID returns Lesson.ID
func (l *Lesson) GetID() int64 {
	return l.ID
}

// GetButtonTitle returns Lesson.Title
func (l *Lesson) GetButtonTitle() string {
	return l.Title
}

// Validate func is needed to validate Lesson object fields before INSERT
//
// NOTE:
// - Title is required
// - Module is required
func (l *Lesson) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Title, validation.Required),
		validation.Field(&l.Module, validation.Required),
	)
}

// ToString converts Lesson object to json string
func (l *Lesson) ToString() string {
	str, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields returns logrus.Fields for logrus logger
//
// NOTE:
// available fields are "id", "title", "module"
func (l *Lesson) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":     l.ID,
		"title":  l.Title,
		"module": l.Module.LogrusFields(),
	}
}
