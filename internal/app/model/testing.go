package model

import "testing"

// TestStudent ...
func TestStudent(y *testing.T) *Student {
	return &Student{
		TelegramID: 99999,
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
	}
}
