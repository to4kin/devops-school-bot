package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		TelegramID: int64(99999),
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
		IsAdmin:    false,
	}
}

func TestSchool(t *testing.T) *School {
	return &School{
		Title:      "2021.2",
		InProgress: true,
		Finished:   false,
	}
}

func TestHomework(t *testing.T) *Homework {
	return &Homework{
		Title: "golang",
	}
}
