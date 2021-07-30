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

func TestSchoolRepository_FindActive(t *testing.T) {
	s := teststore.New()
	testSchool := model.TestSchool(t)

	_, err := s.School().FindActive()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.School().Create(testSchool))

	school, err := s.School().FindActive()
	assert.NoError(t, err)
	assert.NotNil(t, school)
}
