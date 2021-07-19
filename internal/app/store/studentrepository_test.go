package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

var (
	studentTelegramID int    = 99999
	studentFirstName  string = "TestUser"
)

func TestStudentRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("students")

	student, err := s.Student().Create(&model.Student{
		TelegramID: studentTelegramID,
		FirstName:  studentFirstName,
	})
	assert.NoError(t, err)
	assert.NotNil(t, student)
}

func TestStudentRepository_FindBytelegramID(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("students")

	_, err := s.Student().FindByTelegramID(studentTelegramID)
	assert.Error(t, err)

	s.Student().Create(&model.Student{
		TelegramID: studentTelegramID,
		FirstName:  studentFirstName,
	})

	student, err := s.Student().FindByTelegramID(studentTelegramID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}
