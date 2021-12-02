package helper_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestStudentHelper_GetStudent(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	homework := model.TestHomeworkOne(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	replyMessage, replyMarkup, err := hlpr.GetStudent(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestStudentHelper_GetStudentsList(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	student := model.TestStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := hlpr.GetStudentsList(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestStudentHelper_UpdateStudentStatus(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	student := model.TestStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.Equal(t, true, student.Active)
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := hlpr.UpdateStudentStatus(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, false, student.Active)
}

func TestStudentHelper_UpdateStudentType(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	student := model.TestStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.Equal(t, true, student.FullCourse)
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := hlpr.UpdateStudentType(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, false, student.FullCourse)
}
