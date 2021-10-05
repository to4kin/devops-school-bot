package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

func TestHomeworkRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomeworkOne(t)

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))

	assert.NoError(t, s.Homework().Create(h))
	assert.NotNil(t, h)
}

func TestHomeworkRepository_Update(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomeworkOne(t)

	assert.EqualError(t, s.Homework().Update(h), store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))

	assert.NoError(t, s.Homework().Create(h))

	h.Active = false

	assert.NoError(t, s.Homework().Update(h))
	assert.Equal(t, false, h.Active)
}

func TestHomeworkRepository_DeleteByMessageID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h1 := model.TestHomeworkOne(t)
	h2 := model.TestHomeworkTwo(t)

	assert.EqualError(t, s.Homework().DeleteByMessageID(h1.MessageID), store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h1.Student.Account))
	assert.NoError(t, s.School().Create(h1.Student.School))
	assert.NoError(t, s.Student().Create(h1.Student))
	assert.NoError(t, s.Module().Create(h1.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h1.Lesson))
	assert.NoError(t, s.Homework().Create(h1))

	h1.Lesson = h2.Lesson
	assert.NoError(t, s.Lesson().Create(h1.Lesson))
	assert.NoError(t, s.Homework().Create(h1))

	assert.NoError(t, s.Homework().DeleteByMessageID(h1.MessageID))
}

func TestHomeworkRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomeworkOne(t)

	_, err := s.Homework().FindByID(h.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homework, err := s.Homework().FindByID(h.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homework)
}

func TestHomeworkRepository_FindByStudentID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomeworkOne(t)

	_, err := s.Homework().FindByStudentID(h.Student.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homeworks, err := s.Homework().FindByStudentID(h.Student.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homeworks)
}

func TestHomeworkRepository_FindBySchoolID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomeworkOne(t)

	_, err := s.Homework().FindBySchoolID(h.Student.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homeworks, err := s.Homework().FindBySchoolID(h.Student.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homeworks)
}

func TestHomeworkRepository_FindBySchoolIDLessonID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomeworkOne(t)

	_, err := s.Homework().FindBySchoolID(h.Student.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homeworks, err := s.Homework().FindBySchoolIDLessonID(h.Student.School.ID, h.Lesson.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homeworks)
}

func TestHomeworkRepository_FindByStudentIDLessonID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "lesson", "module", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomeworkOne(t)

	_, err := s.Homework().FindByStudentIDLessonID(h.Student.ID, h.Lesson.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	homework, err := s.Homework().FindByStudentIDLessonID(h.Student.ID, h.Lesson.ID)
	assert.NoError(t, err)
	assert.NotNil(t, homework)
}
