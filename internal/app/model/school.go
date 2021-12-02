package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// School object represents a school
type School struct {
	// ID returns a School.ID
	//
	// NOTE: filled in automatically after INSERT to the store
	ID int64 `json:"id,string"`

	// Created returns time.Time
	//
	// NOTE: should be set before INSERT to the store,
	// *field is required
	Created time.Time `json:"created"`

	// Title returns school title
	//
	// NOTE: should be equal to telegram chat title for future search,
	// *field is required
	Title string `json:"title"`

	// ChatID returns chat id where school was started
	//
	// NOTE: should be equal to telegram chat id for future search
	// *field is required
	ChatID int64 `json:"chat_id,string"`

	// Active returns true/false
	//
	// NOTE: false by default if not specified
	Active bool `json:"active,string"`
}

// GetID returns School.ID
func (s *School) GetID() int64 {
	return s.ID
}

// GetStatusIcon returns string depending on active
//
// NOTE:
// 🟢 if active is true
// 🔴 if active is false
func (s *School) GetStatusIcon() string {
	if s.Active {
		return "🟢"
	}

	return "🔴"
}

// GetStatusText returns string depending on active
//
// NOTE:
// 🟢 In Progress if active is true
// 🔴 Stop if active is false
func (s *School) GetStatusText() string {
	if s.Active {
		return "🟢 In Progress"
	}

	return "🔴 Stop"
}

// GetButtonTitle returns composite string depending on active
//
// NOTE: GetStatusIcon() + <space> + Title
func (s *School) GetButtonTitle() string {
	return fmt.Sprintf("%v %v", s.GetStatusIcon(), s.Title)
}

// GetURL returns link to the corresponding telegram chat
//
// NOTE: https://t.me/c/ + ChatID without first 4 characters
func (s *School) GetURL() string {
	return fmt.Sprintf("%v/%v", "https://t.me/c", strconv.FormatInt(s.ChatID, 10)[4:])
}

// Validate func is needed to validate School object fields before INSERT
//
// NOTE:
// - Created is required
// - Title is required
// - ChatID is required
func (s *School) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Created, validation.Required),
		validation.Field(&s.Title, validation.Required),
		validation.Field(&s.ChatID, validation.Required),
	)
}

// ToString converts School object to json string
func (s *School) ToString() string {
	str, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields returns logrus.Fields for logrus logger
//
// NOTE:
// available fields are "id", "title", "chat_id", "active"
func (s *School) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":      s.ID,
		"title":   s.Title,
		"chat_id": s.ChatID,
		"active":  s.Active,
	}
}
