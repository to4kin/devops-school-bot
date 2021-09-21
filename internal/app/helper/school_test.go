package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestSchoolHelper_GetSchoolsList(t *testing.T) {
	store := teststore.New()
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := helper.GetSchoolsList(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_GetSchool(t *testing.T) {
	store := teststore.New()
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := helper.GetSchool(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_StartSchool(t *testing.T) {
	store := teststore.New()
	school := model.TestInactiveSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.Equal(t, false, school.Active)
	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := helper.StartSchool(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, true, school.Active)
}

func TestSchoolHelper_StopSchool(t *testing.T) {
	store := teststore.New()
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.Equal(t, true, school.Active)
	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := helper.StopSchool(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, false, school.Active)
}

func TestSchoolHelper_ReportSchool(t *testing.T) {
	store := teststore.New()
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := helper.ReportSchool(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_FullReportSchool(t *testing.T) {
	store := teststore.New()
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := helper.FullReportSchool(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_GetSchoolHomeworks(t *testing.T) {
	store := teststore.New()
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := helper.GetSchoolHomeworks(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}
