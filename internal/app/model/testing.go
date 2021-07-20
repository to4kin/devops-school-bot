package model

import "testing"

func TestAccount(t *testing.T) *Account {
	return &Account{
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

func TestLesson(t *testing.T) *Leson {
	return &Leson{
		Title: "golang",
	}
}
