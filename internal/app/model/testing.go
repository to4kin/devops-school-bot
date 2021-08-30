package model

import (
	"testing"
	"time"
)

// TestAccount ...
func TestAccount(t *testing.T) *Account {
	return &Account{
		Created:    time.Now(),
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
		Created: time.Now(),
		Title:   "Title",
		ChatID:  int64(99999),
		Active:  true,
	}
}

// TestLesson ...
func TestLesson(t *testing.T) *Lesson {
	return &Lesson{
		Title:  "golang",
		Module: TestModule(t),
	}
}

// TestModule ...
func TestModule(t *testing.T) *Module {
	return &Module{
		Title: "golang",
	}
}

// TestStudent ...
func TestStudent(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     true,
		FullCourse: true,
	}
}

// TestHomework ...
func TestHomework(t *testing.T) *Homework {
	return &Homework{
		Created:   time.Now(),
		Student:   TestStudent(t),
		Lesson:    TestLesson(t),
		MessageID: int64(99999),
		Verify:    true,
	}
}
