package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

func TestLessonRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("lesson")

	s := sqlstore.New(db)
	l := model.TestLesson(t)

	assert.NoError(t, s.Lesson().Create(l))
	assert.NotNil(t, l)
}

func TestLesson_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("lesson")

	s := sqlstore.New(db)
	l := model.TestLesson(t)

	_, err := s.Lesson().FindByID(l.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Lesson().Create(l))

	lesson, err := s.Lesson().FindByID(l.ID)
	assert.NoError(t, err)
	assert.NotNil(t, lesson)
}

func TestLesson_FindByTitle(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("lesson")

	s := sqlstore.New(db)
	l := model.TestLesson(t)

	_, err := s.Lesson().FindByTitle(l.Title)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Lesson().Create(l))

	lesson, err := s.Lesson().FindByTitle(l.Title)
	assert.NoError(t, err)
	assert.NotNil(t, lesson)
}

func TestLesson_FindBySchoolID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homework", "lesson", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomework(t)

	_, err := s.Lesson().FindBySchoolID(h.Student.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	lesson, err := s.Lesson().FindBySchoolID(h.Student.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, lesson)
}
