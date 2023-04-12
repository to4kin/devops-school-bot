package model

import (
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Student object represents a student
type Student struct {
	// ID returns a Student.ID
	//
	// NOTE: filled in automatically after INSERT to the store
	ID int64 `json:"id,string"`

	// Created returns time.Time
	//
	// NOTE: should be set before INSERT to the store,
	// *field is required
	Created time.Time `json:"created"`

	// Account object represents an account
	Account *Account `json:"account"`

	// School object represents a school
	School *School `json:"school"`

	// Active returns true/false
	//
	// NOTE: false by default if not specified
	Active bool `json:"active,string"`

	// FullCourse returns true/false
	//
	// NOTE:
	// - "true" means - "Student"
	// - "false" means - "Listener"
	// false by default if not specified
	FullCourse bool `json:"full_course,string"`
}

// GetID returns Student.ID
func (s *Student) GetID() int64 {
	return s.ID
}

// GetStatusIcon returns string depending on active
//
// NOTE:
// 🟢 if active is true
// 🔴 if active is false
func (s *Student) GetStatusIcon() string {
	if s.Active {
		return "🟢"
	}

	return "🔴"
}

// GetStatusText returns string depending on active
//
// NOTE:
// 🟢Active if active is true
// 🔴Block if active is false
func (s *Student) GetStatusText() string {
	if s.Active {
		return "🟢Active"
	}

	return "🔴Block"
}

// GetButtonTitle returns composite string depending on active
//
// NOTE: GetStatusIcon() + <space> + Account.Username
func (s *Student) GetButtonTitle() string {
	return fmt.Sprintf("%v%v", s.GetStatusIcon(), s.Account.GetFullName())
}

// GetType returns student type
//
// NOTE:
// - type is "Student" if FullCourse is true
// - type is "Listener" if FullCourse is false
func (s *Student) GetType() string {
	if s.FullCourse {
		return "Student"
	}
	return "Listener"
}

// Validate func is needed to validate Student object fields before INSERT
//
// NOTE:
// - Created is required
// - Account is required
// - School is required
func (s *Student) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Created, validation.Required),
		validation.Field(&s.Account, validation.Required),
		validation.Field(&s.School, validation.Required),
	)
}

// ToString converts Student object to json string
func (s *Student) ToString() string {
	str, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields returns logrus.Fields for logrus logger
//
// NOTE:
// available fields are "id", "account", "school", "active", "full_course"
func (s *Student) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":          s.ID,
		"account":     s.Account.LogrusFields(),
		"school":      s.School.LogrusFields(),
		"active":      s.Active,
		"full_course": s.FullCourse,
	}
}
