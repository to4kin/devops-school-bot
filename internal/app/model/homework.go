package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Homework ...
type Homework struct {
	ID        int64    `json:"id"`
	Student   *Student `json:"student"`
	Lesson    *Lesson  `json:"lesson"`
	MessageID int64    `json:"message_id"`
	Verify    bool     `json:"verify"`
}

// Validate ...
func (h *Homework) Validate() error {
	return validation.ValidateStruct(
		h,
		validation.Field(&h.Student, validation.Required),
		validation.Field(&h.Lesson, validation.Required),
		validation.Field(&h.MessageID, validation.Required),
	)
}

// ToString ...
func (h *Homework) ToString() string {
	str, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields ...
func (h *Homework) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":         h.ID,
		"student":    h.Student.LogrusFields(),
		"lesson":     h.Lesson.LogrusFields(),
		"message_id": h.MessageID,
		"verify":     h.Verify,
	}
}
