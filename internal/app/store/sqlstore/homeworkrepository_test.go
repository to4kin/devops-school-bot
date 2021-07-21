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

	s.Account().Create(h.Student.Account)
	s.School().Create(h.Student.School)
	s.Student().Create(h.Student)
	s.Lesson().Create(h.Lesson)

	assert.NoError(t, s.Homework().Create(h))
	assert.NotNil(t, h)
}

func TestHomeworkRepository_FindByStudentIDLessonID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homework", "lesson", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomework(t)

	_, err := s.Homework().FindByStudentIDLessonID(h.Student.ID, h.Lesson.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(h.Student.Account)
	s.School().Create(h.Student.School)
	s.Student().Create(h.Student)
	s.Lesson().Create(h.Lesson)
	s.Homework().Create(h)

	homework, err := s.Homework().FindByStudentIDLessonID(h.Student.ID, h.Lesson.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homework)
}
