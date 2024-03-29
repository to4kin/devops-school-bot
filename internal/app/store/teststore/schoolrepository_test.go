package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestSchoolRepository_Create(t *testing.T) {
	s := teststore.New()
	school := model.TestSchool(t)

	assert.NoError(t, s.School().Create(school))
	assert.NotNil(t, school)

	assert.EqualError(t, s.School().Create(school), store.ErrRecordIsExist.Error())
}

func TestSchoolRepository_Update(t *testing.T) {
	s := teststore.New()
	school := model.TestSchool(t)

	assert.EqualError(t, s.School().Update(school), store.ErrRecordNotFound.Error())
	assert.NoError(t, s.School().Create(school))

	school.Title = "NewTitle"
	school.Active = false

	assert.NoError(t, s.School().Update(school))
	assert.Equal(t, false, school.Active)
}

func TestSchoolRepository_FindAll(t *testing.T) {
	s := teststore.New()
	school := model.TestSchool(t)

	schools, err := s.School().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, schools)
	assert.NoError(t, s.School().Create(school))

	schools, err = s.School().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, schools)
}

func TestSchoolRepository_FindByID(t *testing.T) {
	s := teststore.New()
	testSchool := model.TestSchool(t)

	_, err := s.School().FindByID(testSchool.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.School().Create(testSchool))

	school, err := s.School().FindByID(testSchool.ID)
	assert.NoError(t, err)
	assert.NotNil(t, school)
}

func TestSchoolRepository_FindByTitle(t *testing.T) {
	s := teststore.New()
	testSchool := model.TestSchool(t)

	_, err := s.School().FindByTitle(testSchool.Title)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.School().Create(testSchool))

	school, err := s.School().FindByTitle(testSchool.Title)
	assert.NoError(t, err)
	assert.NotNil(t, school)
}

func TestSchoolRepository_FindByChatID(t *testing.T) {
	s := teststore.New()
	testSchool := model.TestSchool(t)

	_, err := s.School().FindByChatID(testSchool.ChatID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.School().Create(testSchool))

	school, err := s.School().FindByChatID(testSchool.ChatID)
	assert.NoError(t, err)
	assert.NotNil(t, school)
}

func TestSchoolRepository_FindByActive(t *testing.T) {
	s := teststore.New()
	school := model.TestSchool(t)

	schools, err := s.School().FindByActive(school.Active)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, schools)
	assert.NoError(t, s.School().Create(school))

	schools, err = s.School().FindByActive(school.Active)
	assert.NoError(t, err)
	assert.NotNil(t, schools)
	assert.Equal(t, school.Active, schools[0].Active)
}
