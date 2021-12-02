package helper_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestReportHelper_GetUserReport(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	homework := model.TestHomeworkOne(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := hlpr.GetUserReport(homework.Student.Account)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}

func TestReportHelper_GetReport(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	homework := model.TestHomeworkOne(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := hlpr.GetReport(homework.Student.School)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}

func TestReportHelper_GetFullReport(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	homework := model.TestHomeworkOne(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := hlpr.GetFullReport(homework.Student.School)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}

func TestReportHelper_GetLessonsReport(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	homework := model.TestHomeworkOne(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := hlpr.GetLessonsReport(homework.Student.School)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}
