package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestHomeworkHelper_GetHomework(t *testing.T) {
	store := teststore.New()
	homework := model.TestHomework(t)
	callback := model.TestHomeworkCallback(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	replyMessage, replyMarkup, err := helper.GetHomework(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestHomeworkHelper_GetHomeworksList(t *testing.T) {
	store := teststore.New()
	homework := model.TestHomework(t)
	callback := model.TestHomeworkCallback(t)

	assert.NoError(t, store.Account().Create(homework.Student.Account))
	assert.NoError(t, store.School().Create(homework.Student.School))
	assert.NoError(t, store.Student().Create(homework.Student))
	assert.NoError(t, store.Lesson().Create(homework.Lesson))
	assert.NoError(t, store.Homework().Create(homework))

	replyMessage, replyMarkup, err := helper.GetHomeworksList(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}
