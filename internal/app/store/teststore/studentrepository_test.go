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

	assert.EqualError(t, s.Student().Create(testStudent), store.ErrRecordIsExist.Error())
}

func TestStudentRepository_Update(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	assert.EqualError(t, s.Student().Update(testStudent), store.ErrRecordNotFound.Error())
	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))

	assert.NoError(t, s.Student().Create(testStudent))

	testStudent.Active = false

	assert.NoError(t, s.Student().Update(testStudent))
	assert.Equal(t, false, testStudent.Active)
}

func TestStudentRepository_FindAll(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	student, err := s.Student().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, student)
}

func TestStudentRepository_FindByID(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByID(testStudent.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	student, err := s.Student().FindByID(testStudent.ID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, testStudent.ID, student.ID)
}

func TestStudentRepository_FindBySchoolID(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindBySchoolID(testStudent.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	students, err := s.Student().FindBySchoolID(testStudent.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, students)
	assert.Equal(t, testStudent.School.ID, students[0].School.ID)
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
	assert.Equal(t, testStudent.Account.ID, student.Account.ID)
	assert.Equal(t, testStudent.School.ID, student.School.ID)
}

func TestStudentRepository_FindByFullCourseSchoolID(t *testing.T) {
	s := teststore.New()
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByFullCourseSchoolID(testStudent.FullCourse, testStudent.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	students, err := s.Student().FindByFullCourseSchoolID(testStudent.FullCourse, testStudent.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, students)
	assert.Equal(t, testStudent.FullCourse, students[0].FullCourse)
}
