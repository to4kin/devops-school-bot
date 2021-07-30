package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Homework struct {
	ID        int64    `json:"id"`
	Student   *Student `json:"student"`
	Lesson    *Lesson  `json:"lesson"`
	ChatID    int64    `json:"chat_id"`
	MessageID int64    `json:"message_id"`
	Verify    bool     `json:"verify"`
}

func (h *Homework) Validate() error {
	return validation.ValidateStruct(
		h,
		validation.Field(&h.Student, validation.Required),
		validation.Field(&h.Lesson, validation.Required),
		validation.Field(&h.ChatID, validation.Required),
		validation.Field(&h.MessageID, validation.Required),
	)
}

func (h *Homework) ToString() string {
	str, err := json.Marshal(h)
	if err != nil {
		return ""
	}

	return string(str)
}
