package model

import "testing"

func TestStudent(y *testing.T) *Student {
	return &Student{
		TelegramID: int64(99999),
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
	}
}
