package model

import (
	"encoding/json"
	"errors"
)

// Callback ...
type Callback struct {
	Action  string      `json:"action"`
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// ToString ...
func (c *Callback) ToString() string {
	str, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return ""
	}

	return string(str)
}

// UnmarshalJSON ...
func (c *Callback) UnmarshalJSON(b []byte) error {
	type typedObject struct {
		Action  string          `json:"action"`
		Type    string          `json:"type"`
		Content json.RawMessage `json:"content"`
	}

	sr := &typedObject{}
	if err := json.Unmarshal(b, sr); err != nil {
		return err
	}

	var value interface{}
	switch sr.Type {
	case "account":
		value = new(Account)
	case "school":
		value = new(School)
	case "student":
		value = new(Student)
	default:
		return errors.New("unsupported type")
	}

	if err := json.Unmarshal(sr.Content, &value); err != nil {
		return err
	}

	c.Action = sr.Action
	c.Type = sr.Type
	c.Content = value

	return nil
}
