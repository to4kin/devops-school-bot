package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/sqlstore"
)

func TestStudentRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("students")

	s := sqlstore.New(db)
	testStudent := model.TestStudent(t)

	assert.NoError(t, s.Student().Create(testStudent))
	assert.NotNil(t, testStudent)
}

func TestStudentRepository_FindBytelegramID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("students")

	s := sqlstore.New(db)
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByTelegramID(testStudent.TelegramID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Student().Create(testStudent)

	student, err := s.Student().FindByTelegramID(testStudent.TelegramID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}
