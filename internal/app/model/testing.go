package model

import "testing"

func TestAccount(t *testing.T) *Account {
	return &Account{
		TelegramID: int64(99999),
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
		Superuser:  false,
	}
}

func TestSchool(t *testing.T) *School {
	return &School{
		Title:    "2021.2",
		Active:   true,
		Finished: false,
	}
}

func TestLesson(t *testing.T) *Lesson {
	return &Lesson{
		Title: "golang",
	}
}

func TestStudent(t *testing.T) *Student {
	return &Student{
		Account: TestAccount(t),
		School:  TestSchool(t),
		Active:  true,
	}
}

func TestHomework(t *testing.T) *Homework {
	return &Homework{
		Student:   TestStudent(t),
		Lesson:    TestLesson(t),
		ChatID:    int64(99999),
		MessageID: int64(99999),
		Verify:    true,
	}
}
