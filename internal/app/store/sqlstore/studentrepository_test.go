package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

func TestStudentRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("student", "school", "account")

	s := sqlstore.New(db)
	testStudent := model.TestStudent(t)

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))

	assert.NoError(t, s.Student().Create(testStudent))
	assert.NotNil(t, testStudent)
}

func TestStudentRepository_Update(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("student", "school", "account")

	s := sqlstore.New(db)
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
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("student", "school", "account")

	s := sqlstore.New(db)
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
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("student", "school", "account")

	s := sqlstore.New(db)
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindByID(testStudent.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	student, err := s.Student().FindByID(testStudent.ID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}

func TestStudentRepository_FindBySchoolID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("student", "school", "account")

	s := sqlstore.New(db)
	testStudent := model.TestStudent(t)

	_, err := s.Student().FindBySchoolID(testStudent.School.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(testStudent.Account))
	assert.NoError(t, s.School().Create(testStudent.School))
	assert.NoError(t, s.Student().Create(testStudent))

	student, err := s.Student().FindBySchoolID(testStudent.School.ID)
	assert.NoError(t, err)
	assert.NotNil(t, student)
}

func TestStudentRepository_FindByAccountIDSchoolID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("student", "school", "account")

	s := sqlstore.New(db)
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
