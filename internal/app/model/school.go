package model

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// School ...
type School struct {
	ID       int64     `json:"id"`
	Created  time.Time `json:"created"`
	Title    string    `json:"title"`
	ChatID   int64     `json:"chat_id"`
	Finished bool      `json:"finished"`
}

// Validate ...
func (s *School) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Created, validation.Required),
		validation.Field(&s.Title, validation.Required),
		validation.Field(&s.ChatID, validation.Required),
	)
}

// ToString ...
func (s *School) ToString() string {
	str, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields ...
func (s *School) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":       s.ID,
		"title":    s.Title,
		"chat_id":  s.ChatID,
		"finished": s.Finished,
	}
}
