package model

import "testing"

// TestAccount ...
func TestAccount(t *testing.T) *Account {
	return &Account{
		TelegramID: int64(99999),
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
		Superuser:  false,
	}
}

// TestSchool ...
func TestSchool(t *testing.T) *School {
	return &School{
		Title:    "DevOps School 2021.2",
		ChatID:   int64(99999),
		Active:   true,
		Finished: false,
	}
}

// TestLesson ...
func TestLesson(t *testing.T) *Lesson {
	return &Lesson{
		Title: "golang",
	}
}

// TestStudent ...
func TestStudent(t *testing.T) *Student {
	return &Student{
		Account: TestAccount(t),
		School:  TestSchool(t),
		Active:  true,
	}
}

// TestHomework ...
func TestHomework(t *testing.T) *Homework {
	return &Homework{
		Student:   TestStudent(t),
		Lesson:    TestLesson(t),
		MessageID: int64(99999),
		Verify:    true,
	}
}
