package model

import (
	"encoding/json"
)

// Callback ...
type Callback struct {
	Type string `json:"type"`
	ID   int64  `json:"id,string"`
}

// ToString ...
func (c *Callback) ToString() string {
	str, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(str)
}

// Unmarshal ...
func (c *Callback) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, c); err != nil {
		return err
	}

	return nil
}
