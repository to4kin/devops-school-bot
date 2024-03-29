package helper_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestSchoolHelper_GetSchoolsList(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := hlpr.GetSchoolsList(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_GetSchool(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := hlpr.GetSchool(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_StartSchool(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	school := model.TestInactiveSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.Equal(t, false, school.Active)
	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := hlpr.UpdateSchoolStatus(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, true, school.Active)
}

func TestSchoolHelper_GetSchoolReport(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := hlpr.GetSchoolReport(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_GetSchoolFullReport(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := hlpr.GetSchoolFullReport(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestSchoolHelper_GetSchoolHomeworks(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	school := model.TestSchool(t)
	callback := model.TestSchoolCallback(t)

	assert.NoError(t, store.School().Create(school))
	replyMessage, replyMarkup, err := hlpr.GetSchoolHomeworks(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}
