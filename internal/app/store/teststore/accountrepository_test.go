package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/teststore"
)

func TestAccountRepository_Create(t *testing.T) {
	s := teststore.New()
	a := model.TestAccount(t)

	assert.NoError(t, s.Account().Create(a))
	assert.NotNil(t, a)
}

func TestAccountRepository_FindByTelegramID(t *testing.T) {
	s := teststore.New()
	a := model.TestAccount(t)

	_, err := s.Account().FindByTelegramID(a.TelegramID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(a)

	account, err := s.Account().FindByTelegramID(a.TelegramID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
}
