package model

import (
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// School ...
type School struct {
	ID      int64     `json:"id,string"`
	Created time.Time `json:"created"`
	Title   string    `json:"title"`
	ChatID  int64     `json:"chat_id,string"`
	Active  bool      `json:"active,string"`
}

// GetID ...
func (s *School) GetID() int64 {
	return s.ID
}

// GetStatusText ...
func (s *School) GetStatusText() string {
	if s.Active {
		return "🟢"
	}

	return "🔴"
}

// GetButtonTitle ...
func (s *School) GetButtonTitle() string {
	return fmt.Sprintf("%v %v", s.GetStatusText(), s.Title)
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
		"id":      s.ID,
		"title":   s.Title,
		"chat_id": s.ChatID,
		"active":  s.Active,
	}
}
