package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestStudentHelper_GetStudent(t *testing.T) {
	store := teststore.New()
	student := model.TestStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := helper.GetStudent(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestStudentHelper_GetStudentsList(t *testing.T) {
	store := teststore.New()
	student := model.TestStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := helper.GetStudentsList(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestStudentHelper_BlockStudent(t *testing.T) {
	store := teststore.New()
	student := model.TestStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.Equal(t, true, student.Active)
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := helper.BlockStudent(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, false, student.Active)
}

func TestStudentHelper_UnblockStudent(t *testing.T) {
	store := teststore.New()
	student := model.TestInactiveStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.Equal(t, false, student.Active)
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := helper.UnblockStudent(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, true, student.Active)
}

func TestStudentHelper_SetStudent(t *testing.T) {
	store := teststore.New()
	student := model.TestListener(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.Equal(t, false, student.FullCourse)
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := helper.SetStudent(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, true, student.FullCourse)
}

func TestStudentHelper_SetListener(t *testing.T) {
	store := teststore.New()
	student := model.TestStudent(t)
	callback := model.TestStudentCallback(t)

	assert.NoError(t, store.Account().Create(student.Account))
	assert.NoError(t, store.School().Create(student.School))
	assert.Equal(t, true, student.FullCourse)
	assert.NoError(t, store.Student().Create(student))

	replyMessage, replyMarkup, err := helper.SetListener(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
	assert.Equal(t, false, student.FullCourse)
}
