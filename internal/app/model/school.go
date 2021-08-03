package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// School ...
type School struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	ChatID   int64  `json:"chat_id"`
	Active   bool   `json:"active"`
	Finished bool   `json:"finished"`
}

// Validate ...
func (s *School) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Title, validation.Required),
		validation.Field(&s.ChatID, validation.Required),
		validation.Field(&s.Active, validation.By(notEqual(s.Finished))),
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
		"active":   s.Active,
		"finished": s.Finished,
	}
}
