package model

import (
	"testing"
	"time"

	"gopkg.in/tucnak/telebot.v3"
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

// TestAdminAccount ...
func TestAdminAccount(t *testing.T) *Account {
	return &Account{
		Created:    time.Now(),
		TelegramID: int64(99999),
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
		Superuser:  true,
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

// TestInactiveSchool ...
func TestInactiveSchool(t *testing.T) *School {
	return &School{
		Created: time.Now(),
		Title:   "Title",
		ChatID:  int64(99999),
		Active:  false,
	}
}

// TestLessonOne ...
func TestLessonOne(t *testing.T) *Lesson {
	return &Lesson{
		Title:  "golang1",
		Module: TestModule(t),
	}
}

// TestLessonTwo ...
func TestLessonTwo(t *testing.T) *Lesson {
	return &Lesson{
		Title:  "golang2",
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

// TestInactiveStudent ...
func TestInactiveStudent(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     false,
		FullCourse: true,
	}
}

// TestListener ...
func TestListener(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     true,
		FullCourse: false,
	}
}

// TestInactiveListener ...
func TestInactiveListener(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     false,
		FullCourse: false,
	}
}

// TestHomeworkOne ...
func TestHomeworkOne(t *testing.T) *Homework {
	return &Homework{
		Created:   time.Now(),
		Student:   TestStudent(t),
		Lesson:    TestLessonOne(t),
		MessageID: int64(99999),
		Verify:    true,
		Active:    true,
	}
}

// TestHomeworkTwo ...
func TestHomeworkTwo(t *testing.T) *Homework {
	return &Homework{
		Created:   time.Now(),
		Student:   TestStudent(t),
		Lesson:    TestLessonTwo(t),
		MessageID: int64(99999),
		Verify:    true,
		Active:    true,
	}
}

// TestAccountCallback ...
func TestAccountCallback(t *testing.T) *Callback {
	return &Callback{
		ID:          int64(1),
		Type:        "account",
		Command:     "get",
		ListCommand: "accounts_list",
	}
}

// TestStudentCallback ...
func TestStudentCallback(t *testing.T) *Callback {
	return &Callback{
		ID:          int64(1),
		Type:        "student",
		Command:     "get",
		ListCommand: "students_list",
	}
}

// TestSchoolCallback ...
func TestSchoolCallback(t *testing.T) *Callback {
	return &Callback{
		ID:          int64(1),
		Type:        "school",
		Command:     "get",
		ListCommand: "schools_list",
	}
}

// TestHomeworkCallback ...
func TestHomeworkCallback(t *testing.T) *Callback {
	return &Callback{
		ID:          int64(1),
		Type:        "homework",
		Command:     "get",
		ListCommand: "homeworks_list",
	}
}

// TestTelebotUser ...
func TestTelebotUser(t *testing.T) *telebot.User {
	return &telebot.User{
		ID:        int64(99999),
		FirstName: "FirstName",
		LastName:  "LastName",
		Username:  "Username",
	}
}
