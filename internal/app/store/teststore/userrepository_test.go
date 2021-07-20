package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/teststore"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindBytelegramID(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)

	_, err := s.User().FindByTelegramID(u.TelegramID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(u)

	user, err := s.User().FindByTelegramID(u.TelegramID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
