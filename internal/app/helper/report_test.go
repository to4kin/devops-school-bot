package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestReportHelper_GetUserReport(t *testing.T) {
	store := teststore.New()
	homework := model.TestHomework(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := helper.GetUserReport(store, homework.Student.Account, homework.Student.School)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}

func TestReportHelper_GetReport(t *testing.T) {
	store := teststore.New()
	homework := model.TestHomework(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := helper.GetReport(store, homework.Student.School)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}

func TestReportHelper_GetFullReport(t *testing.T) {
	store := teststore.New()
	homework := model.TestHomework(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := helper.GetFullReport(store, homework.Student.School)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}

func TestReportHelper_GetLessonsReport(t *testing.T) {
	store := teststore.New()
	homework := model.TestHomework(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	reportMessage, err := helper.GetLessonsReport(store, homework.Student.School)
	assert.NoError(t, err)
	assert.NotEmpty(t, reportMessage)
}