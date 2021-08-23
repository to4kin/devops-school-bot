package model

import (
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Student ...
type Student struct {
	ID         int64     `json:"id,string"`
	Created    time.Time `json:"created"`
	Account    *Account  `json:"account"`
	School     *School   `json:"school"`
	Active     bool      `json:"active,string"`
	FullCourse bool      `json:"full_course,string"`
}

// GetID ...
func (s *Student) GetID() int64 {
	return s.ID
}

// GetStatusText ...
func (s *Student) GetStatusText() string {
	if s.Active {
		return "🟢"
	}

	return "🔴"
}

// GetButtonTitle ...
func (s *Student) GetButtonTitle() string {
	return fmt.Sprintf("%v @%v", s.GetStatusText(), s.Account.Username)
}

// GetType ...
func (s *Student) GetType() string {
	if s.FullCourse {
		return "Student"
	}
	return "Listener"
}

// Validate ...
func (s *Student) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Created, validation.Required),
		validation.Field(&s.Account, validation.Required),
		validation.Field(&s.School, validation.Required),
	)
}

// ToString ...
func (s *Student) ToString() string {
	str, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields ...
func (s *Student) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":          s.ID,
		"account":     s.Account.LogrusFields(),
		"school":      s.School.LogrusFields(),
		"active":      s.Active,
		"full_course": s.FullCourse,
	}
}
