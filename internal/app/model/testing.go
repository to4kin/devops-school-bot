package model

import (
	"testing"
	"time"

	"gopkg.in/tucnak/telebot.v3"
)

// TestAccount for testing
// NOTE:
//
//	Created:    time.Now(),
//	TelegramID: int64(99999),
//	FirstName:  "FirstName",
//	LastName:   "LastName",
//	Username:   "Username",
//	Superuser:  false,
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

// TestAdminAccount for testing
// NOTE:
//
//	Created:    time.Now(),
//	TelegramID: int64(99999),
//	FirstName:  "FirstName",
//	LastName:   "LastName",
//	Username:   "Username",
//	Superuser:  true,
func TestAdminAccount(t *testing.T) *Account {
	return &Account{
		Created:    time.Now(),
		TelegramID: int64(66666),
		FirstName:  "FirstName",
		LastName:   "LastName",
		Username:   "Username",
		Superuser:  true,
	}
}

// TestSchool for testing
// NOTE:
//
//	Created: time.Now(),
//	Title:   "Title",
//	ChatID:  int64(99999),
//	Active:  true,
func TestSchool(t *testing.T) *School {
	return &School{
		Created: time.Now(),
		Title:   "Title",
		ChatID:  int64(99999),
		Active:  true,
	}
}

// TestInactiveSchool for testing
// NOTE:
//
//	Created: time.Now(),
//	Title:   "Title",
//	ChatID:  int64(99999),
//	Active:  false,
func TestInactiveSchool(t *testing.T) *School {
	return &School{
		Created: time.Now(),
		Title:   "Title",
		ChatID:  int64(99999),
		Active:  false,
	}
}

// TestLessonOne for testing
// NOTE:
//
//	Title:  "golang1",
//	Module: TestModule(t),
func TestLessonOne(t *testing.T) *Lesson {
	return &Lesson{
		Title:  "golang1",
		Module: TestModule(t),
	}
}

// TestLessonTwo for testing
// NOTE:
//
//	Title:  "golang2",
//	Module: TestModule(t),
func TestLessonTwo(t *testing.T) *Lesson {
	return &Lesson{
		Title:  "golang2",
		Module: TestModule(t),
	}
}

// TestModule for testing
// NOTE:
//
//	Title: "golang",
func TestModule(t *testing.T) *Module {
	return &Module{
		Title: "golang",
	}
}

// TestStudent for testing
// NOTE:
//
//	Created:    time.Now(),
//	Account:    TestAccount(t),
//	School:     TestSchool(t),
//	Active:     true,
//	FullCourse: true,
func TestStudent(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     true,
		FullCourse: true,
	}
}

// TestInactiveStudent for testing
// NOTE:
//
//	Created:    time.Now(),
//	Account:    TestAccount(t),
//	School:     TestSchool(t),
//	Active:     false,
//	FullCourse: true,
func TestInactiveStudent(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     false,
		FullCourse: true,
	}
}

// TestListener for testing
// NOTE:
//
//	Created:    time.Now(),
//	Account:    TestAccount(t),
//	School:     TestSchool(t),
//	Active:     true,
//	FullCourse: false,
func TestListener(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     true,
		FullCourse: false,
	}
}

// TestInactiveListener for testing
// NOTE:
//
//	Created:    time.Now(),
//	Account:    TestAccount(t),
//	School:     TestSchool(t),
//	Active:     false,
//	FullCourse: false,
func TestInactiveListener(t *testing.T) *Student {
	return &Student{
		Created:    time.Now(),
		Account:    TestAccount(t),
		School:     TestSchool(t),
		Active:     false,
		FullCourse: false,
	}
}

// TestHomeworkOne for testing
// NOTE:
//
//	Created:   time.Now(),
//	Student:   TestStudent(t),
//	Lesson:    TestLessonOne(t),
//	MessageID: int64(99999),
//	Verify:    true,
//	Active:    true,
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

// TestHomeworkTwo for testing
// NOTE:
//
//	Created:   time.Now(),
//	Student:   TestStudent(t),
//	Lesson:    TestLessonTwo(t),
//	MessageID: int64(99999),
//	Verify:    true,
//	Active:    true,
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

// TestAccountCallback for testing
// NOTE:
//
//	Created:     time.Now(),
//	Type:        "account",
//	TypeID:      int64(1),
//	Command:     "get",
//	ListCommand: "accounts_list",
func TestAccountCallback(t *testing.T) *Callback {
	return &Callback{
		Created:     time.Now(),
		Type:        "account",
		TypeID:      int64(1),
		Command:     "get",
		ListCommand: "accounts_list",
	}
}

// TestStudentCallback for testing
// NOTE:
//
//	Created:     time.Now(),
//	Type:        "student",
//	TypeID:      int64(1),
//	Command:     "get",
//	ListCommand: "students_list",
func TestStudentCallback(t *testing.T) *Callback {
	return &Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      int64(1),
		Command:     "get",
		ListCommand: "students_list",
	}
}

// TestSchoolCallback for testing
// NOTE:
//
//	Created:     time.Now(),
//	Type:        "school",
//	TypeID:      int64(1),
//	Command:     "get",
//	ListCommand: "schools_list",
func TestSchoolCallback(t *testing.T) *Callback {
	return &Callback{
		Created:     time.Now(),
		Type:        "school",
		TypeID:      int64(1),
		Command:     "get",
		ListCommand: "schools_list",
	}
}

// TestHomeworkCallback for testing
// NOTE:
//
//	Created:     time.Now(),
//	Type:        "homework",
//	TypeID:      int64(1),
//	Command:     "get",
//	ListCommand: "homeworks_list",
func TestHomeworkCallback(t *testing.T) *Callback {
	return &Callback{
		Created:     time.Now(),
		Type:        "homework",
		TypeID:      int64(1),
		Command:     "get",
		ListCommand: "homeworks_list",
	}
}

// TestTelebotUser for testing
// NOTE:
//
//	ID:        int64(99999),
//	FirstName: "FirstName",
//	LastName:  "LastName",
//	Username:  "Username",
func TestTelebotUser(t *testing.T) *telebot.User {
	return &telebot.User{
		ID:        int64(99999),
		FirstName: "FirstName",
		LastName:  "LastName",
		Username:  "Username",
	}
}
