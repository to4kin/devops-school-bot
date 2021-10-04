package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestHomeworkRepository_Create(t *testing.T) {
	s := teststore.New()
	h := model.TestHomework(t)

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))

	assert.NoError(t, s.Homework().Create(h))
	assert.NotNil(t, h)
	assert.EqualError(t, s.Homework().Create(h), store.ErrRecordIsExist.Error())
}

func TestHomeworkRepository_Update(t *testing.T) {
	s := teststore.New()
	h := model.TestHomework(t)

	assert.EqualError(t, s.Homework().Update(h), store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	h.Active = false

	assert.NoError(t, s.Homework().Update(h))
	assert.Equal(t, false, h.Active)
}

func TestHomeworkRepository_FindByID(t *testing.T) {
	s := teststore.New()
	h := model.TestHomework(t)

	_, err := s.Homework().FindByID(h.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homework, err := s.Homework().FindByID(h.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homework)
}

func TestHomeworkRepository_FindByStudentID(t *testing.T) {
	s := teststore.New()
	h := model.TestHomework(t)

	_, err := s.Homework().FindByStudentID(h.Student.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homeworks, err := s.Homework().FindByStudentID(h.Student.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homeworks)
}

func TestHomeworkRepository_FindBySchoolID(t *testing.T) {
	s := teststore.New()
	h := model.TestHomework(t)

	_, err := s.Homework().FindBySchoolID(h.Student.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homeworks, err := s.Homework().FindBySchoolID(h.Student.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homeworks)
}

func TestHomeworkRepository_FindBySchoolIDLessonID(t *testing.T) {
	s := teststore.New()
	h := model.TestHomework(t)

	_, err := s.Homework().FindBySchoolID(h.Student.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homeworks, err := s.Homework().FindBySchoolIDLessonID(h.Student.School.ID, h.Lesson.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homeworks)
}

func TestHomeworkRepository_FindByStudentIDLessonID(t *testing.T) {
	s := teststore.New()
	h := model.TestHomework(t)

	_, err := s.Homework().FindByStudentIDLessonID(h.Student.ID, h.Lesson.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homework, err := s.Homework().FindByStudentIDLessonID(h.Student.ID, h.Lesson.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homework)
}
