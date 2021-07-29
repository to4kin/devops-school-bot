package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/sqlstore"
)

func TestHomeworkRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homework", "lesson", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomework(t)

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))

	assert.NoError(t, s.Homework().Create(h))
	assert.NotNil(t, h)
}

func TestHomeworkRepository_FindByStudentID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homework", "lesson", "student", "school", "account")

	s := sqlstore.New(db)
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
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homework", "lesson", "student", "school", "account")

	s := sqlstore.New(db)
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

func TestHomeworkRepository_FindByStudentIDLessonID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homework", "lesson", "student", "school", "account")

	s := sqlstore.New(db)
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
