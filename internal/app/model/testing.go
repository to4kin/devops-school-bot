package model

import "testing"

func TestStudent(t *testing.T) *Student {
	return &Student{
		TelegramID: int64(99999),
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
	}
}

func TestSchool(t *testing.T) *School {
	return &School{
		Title:      "2021.2",
		InProgress: true,
		Finished:   false,
	}
}
