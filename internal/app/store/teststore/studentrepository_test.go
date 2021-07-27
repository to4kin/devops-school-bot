package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/teststore"
)

func TestStudentRepository_Create(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))

	assert.NoError(t, s.Student().Create(testStudent))
	assert.NotNil(t, testStudent)
}

func TestStudentRepository_FindByAccountSchool(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByAccountSchool(testStudent.Account, testStudent.School)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	student, err := s.Student().FindByAccountSchool(testStudent.Account, testStudent.School)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}
