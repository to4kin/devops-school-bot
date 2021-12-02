package model

import (
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Homework object represents a homework
type Homework struct {
	// ID returns A Homework.ID
	//
	// NOTE: filled in automatically after INSERT to the store
	ID int64 `json:"id"`

	// Created returns time.Time
	//
	// NOTE: should be set before INSERT to the store
	// *field is required
	Created time.Time `json:"created"`

	// Student object
	//
	// NOTE:
	// *object is required
	Student *Student `json:"student"`

	// Lesson object
	//
	// NOTE:
	// *object is required
	Lesson *Lesson `json:"lesson"`

	// MessageID returns telegram message id
	//
	// NOTE: should be equal to real telegram message id for future use
	MessageID int64 `json:"message_id"`

	// Verify returns true/false
	//
	// NOTE:
	// - "true" means - homework is verified
	// - "false" meanss - homework is NOT verified
	// false by default if not specified
	Verify bool `json:"verify"`

	// Active returns true/false
	//
	// NOTE: false by default if not specified
	Active bool `json:"active"`
}

// GetID returns Homework.ID
func (h *Homework) GetID() int64 {
	return h.ID
}

// GetStatusIcon returns string depending on active
//
// NOTE:
// 🟢 if active is true
// 🔴 if active is false
func (h *Homework) GetStatusIcon() string {
	if h.Active {
		return "🟢"
	}

	return "🔴"
}

// GetStatusText returns string depending on active
//
// NOTE:
// 🟢 if active is true
// 🔴 if active is false
func (h *Homework) GetStatusText() string {
	if h.Active {
		return "🟢 Enable"
	}

	return "🔴 Disable"
}

// GetButtonTitle returns composite string depending on active
//
// NOTE: GetStatusIcon() + <sapce> + Lesson.Title
func (h *Homework) GetButtonTitle() string {
	return fmt.Sprintf("%v %v", h.GetStatusIcon(), h.Lesson.Title)
}

// GetURL reuturns link to the corresponding telegram message
//
// NOTE: School.GetURL() + / + Homework.MessageID
func (h *Homework) GetURL() string {
	return fmt.Sprintf("%v/%v", h.Student.School.GetURL(), h.MessageID)
}

// Validate func is needed to validate Homework object fields before INSERT
//
// NOTE:
// - Created is required
// - Student is required
// - Lesson is required
// - MessageID is required
func (h *Homework) Validate() error {
	return validation.ValidateStruct(
		h,
		validation.Field(&h.Created, validation.Required),
		validation.Field(&h.Student, validation.Required),
		validation.Field(&h.Lesson, validation.Required),
		validation.Field(&h.MessageID, validation.Required),
	)
}

// ToString converts Homework object to json string
func (h *Homework) ToString() string {
	str, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields returns logrus.Fields for logrus logger
//
// NOTE:
// available fields are "id", "student", "lesson", "message_id", "verify", "active"
func (h *Homework) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":         h.ID,
		"student":    h.Student.LogrusFields(),
		"lesson":     h.Lesson.LogrusFields(),
		"message_id": h.MessageID,
		"verify":     h.Verify,
		"active":     h.Active,
	}
}
