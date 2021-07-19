package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

func TestStudentRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("students")

	student, err := s.Student().Create(model.TestStudent(t))
	assert.NoError(t, err)
	assert.NotNil(t, student)
}

func TestStudentRepository_FindBytelegramID(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("students")

	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByTelegramID(testStudent.TelegramID)
	assert.Error(t, err)

	s.Student().Create(testStudent)

	student, err := s.Student().FindByTelegramID(testStudent.TelegramID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}
