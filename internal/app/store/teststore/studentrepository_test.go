package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestStudentRepository_Create(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))

	assert.NoError(t, s.Student().Create(testStudent))
	assert.NotNil(t, testStudent)
}

func TestStudentRepository_FindByAccountIDSchoolID(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByAccountIDSchoolID(testStudent.Account.ID, testStudent.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	student, err := s.Student().FindByAccountIDSchoolID(testStudent.Account.ID, testStudent.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}
