package helper_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/teststore"
)

func TestUserHelper_GetUsersList(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	account := model.TestAccount(t)
	callback := model.TestAccountCallback(t)

	assert.NoError(t, store.Account().Create(account))
	replyMessage, replyMarkup, err := hlpr.GetUsersList(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestUserHelper_GetUser(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	account := model.TestAccount(t)
	user := model.TestTelebotUser(t)
	callback := model.TestAccountCallback(t)

	assert.NoError(t, store.Account().Create(account))
	replyMessage, replyMarkup, err := hlpr.GetUser(callback, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotEmpty(t, replyMarkup)
}

func TestUserHelper_UpdateUser(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	account := model.TestAccount(t)
	user := model.TestTelebotUser(t)
	callback := model.TestAccountCallback(t)

	assert.NoError(t, store.Account().Create(account))
	replyMessage, replyMarkup, err := hlpr.UpdateUser(callback, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)
}

func TestUserHelper_SetSuperuser(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	account := model.TestAccount(t)
	callback := model.TestAccountCallback(t)

	assert.Equal(t, false, account.Superuser)
	assert.NoError(t, store.Account().Create(account))

	replyMessage, replyMarkup, err := hlpr.SetSuperuser(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)

	assert.Equal(t, true, account.Superuser)
}

func TestUserHelper_UnsetSuperuser(t *testing.T) {
	store := teststore.New()
	hlpr := helper.NewHelper(store, logrus.New())
	account := model.TestAdminAccount(t)
	callback := model.TestAccountCallback(t)

	assert.Equal(t, true, account.Superuser)
	assert.NoError(t, store.Account().Create(account))

	replyMessage, replyMarkup, err := hlpr.UnsetSuperuser(callback)
	assert.NoError(t, err)
	assert.NotEmpty(t, replyMessage)
	assert.NotNil(t, replyMarkup)

	assert.Equal(t, false, account.Superuser)
}
