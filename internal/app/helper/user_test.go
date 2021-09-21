package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestUserHelper_GetUsersList(t *testing.T) {
	store := teststore.New()
	account := model.TestAccount(t)
	callback := model.TestAccountCallback(t)

	assert.NoError(t, store.Account().Create(account))
	replyMessage, replyMarkup, err := helper.GetUsersList(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestUserHelper_GetUser(t *testing.T) {
	store := teststore.New()
	account := model.TestAccount(t)
	user := model.TestTelebotUser(t)
	callback := model.TestAccountCallback(t)

	assert.NoError(t, store.Account().Create(account))
	replyMessage, replyMarkup, err := helper.GetUser(store, callback, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotEmpty(t, replyMarkup)
}

func TestUserHelper_UpdateUser(t *testing.T) {
	store := teststore.New()
	account := model.TestAccount(t)
	user := model.TestTelebotUser(t)
	callback := model.TestAccountCallback(t)

	assert.NoError(t, store.Account().Create(account))
	replyMessage, replyMarkup, err := helper.UpdateUser(store, callback, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestUserHelper_SetSuperuser(t *testing.T) {
	store := teststore.New()
	account := model.TestAccount(t)
	callback := model.TestAccountCallback(t)

	assert.Equal(t, false, account.Superuser)
	assert.NoError(t, store.Account().Create(account))

	replyMessage, replyMarkup, err := helper.SetSuperuser(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)

	assert.Equal(t, true, account.Superuser)
}

func TestUserHelper_UnsetSuperuser(t *testing.T) {
	store := teststore.New()
	account := model.TestAdminAccount(t)
	callback := model.TestAccountCallback(t)

	assert.Equal(t, true, account.Superuser)
	assert.NoError(t, store.Account().Create(account))

	replyMessage, replyMarkup, err := helper.UnsetSuperuser(store, callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)

	assert.Equal(t, false, account.Superuser)
}
