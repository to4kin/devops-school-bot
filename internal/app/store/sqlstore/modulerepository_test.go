package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

func TestModuleRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("module")

	s := sqlstore.New(db)
	m := model.TestModule(t)

	assert.NoError(t, s.Module().Create(m))
	assert.NotNil(t, m)
}

func TestModule_FindAll(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("module")

	s := sqlstore.New(db)
	m := model.TestModule(t)

	_, err := s.Module().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Module().Create(m))

	modules, err := s.Module().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, modules)
}

func TestModule_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("module")

	s := sqlstore.New(db)
	m := model.TestModule(t)

	_, err := s.Module().FindByID(m.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Module().Create(m))

	module, err := s.Module().FindByID(m.ID)
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.Equal(t, m.ID, module.ID)
}

func TestModule_FindByTitle(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("module")

	s := sqlstore.New(db)
	m := model.TestModule(t)

	_, err := s.Module().FindByTitle(m.Title)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Module().Create(m))

	module, err := s.Module().FindByTitle(m.Title)
	assert.NoError(t, err)
	assert.NotNil(t, module)
	assert.Equal(t, m.Title, module.Title)
}

func TestModule_FindBySchoolID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("homework", "module", "lesson", "student", "school", "account")

	s := sqlstore.New(db)
	h := model.TestHomework(t)

	_, err := s.Module().FindBySchoolID(h.Student.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(h.Student.Account))
	assert.NoError(t, s.School().Create(h.Student.School))
	assert.NoError(t, s.Student().Create(h.Student))
	assert.NoError(t, s.Module().Create(h.Lesson.Module))
	assert.NoError(t, s.Lesson().Create(h.Lesson))
	assert.NoError(t, s.Homework().Create(h))

	modules, err := s.Module().FindBySchoolID(h.Student.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, modules)
}
