package model

import (
	"encoding/json"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
)

// Callback represents a telebot.Callback().Data() for inline messages
type Callback struct {
	// ID returns a Callback.ID
	//
	// NOTE: filled in automatically after INSERT to the store
	ID int64 `json:"id"`

	// Created returns time.Time when the record was created
	//
	// NOTE: should be set before INSERT to the store,
	// *field is required
	Created time.Time `json:"created"`

	// Type returns type of callback element
	//
	// NOTE: field can be one of "Account", "Homework", "Lesson", "Module"
	// "School", "Student"
	Type string `json:"type"`

	// TypeID returns id of type element
	TypeID int64 `json:"type_id"`

	// Command for button
	Command string `json:"command"`

	// Command for elements in list
	ListCommand string `json:"list_command"`
}

// GetStringID converts int64 ID to a string with base = 10
func (c *Callback) GetStringID() string {
	return strconv.FormatInt(c.ID, 10)
}

// Validate func is needed to validate Callback object before INSERT
//
// NOTE:
// - Created is required
// - Type is required
// - TypeID is required
// - Command is required
// - ListCommand is required
func (c *Callback) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Created, validation.Required),
		validation.Field(&c.Type, validation.Required),
		validation.Field(&c.TypeID, validation.Required),
		validation.Field(&c.Command, validation.Required),
		validation.Field(&c.ListCommand, validation.Required),
	)
}

// ToString converts Callback object to json string
func (c *Callback) ToString() string {
	str, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// LogrusFields returns logrus.Fields for logrus logger
//
// NOTE:
// available fields are "id", "created", "type", "type_id", "command", "list_command"
func (c *Callback) LogrusFields() logrus.Fields {
	return logrus.Fields{
		"id":           c.ID,
		"created":      c.Created,
		"type":         c.Type,
		"type_id":      c.TypeID,
		"command":      c.Command,
		"list_command": c.ListCommand,
	}
}
