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

	s.Account().Create(testStudent.Account)
	s.School().Create(testStudent.School)

	assert.NoError(t, s.Student().Create(testStudent))
	assert.NotNil(t, testStudent)
}

func TestStudentRepository_FindByAccountIDSchoolID(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByAccountIDSchoolID(testStudent.Account.ID, testStudent.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(testStudent.Account)
	s.School().Create(testStudent.School)
	s.Student().Create(testStudent)

	student, err := s.Student().FindByAccountIDSchoolID(testStudent.Account.ID, testStudent.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}
