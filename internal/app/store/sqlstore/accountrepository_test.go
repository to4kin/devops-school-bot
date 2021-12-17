package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

func TestAccountRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	assert.Error(t, s.Account().Create(&model.Account{}))
	assert.NoError(t, s.Account().Create(a))
	assert.NotNil(t, a)
}

func TestAccountRepository_Update(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	assert.Error(t, s.Account().Update(&model.Account{}))
	assert.EqualError(t, s.Account().Update(a), store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(a))
	assert.NotNil(t, a)

	a.FirstName += "New"
	a.LastName += "New"
	a.Username += "New"
	a.Superuser = !a.Superuser

	assert.NoError(t, s.Account().Update(a))

	account, err := s.Account().FindByID(a.ID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, a.Superuser, account.Superuser)
}

func TestAccountRepository_FindAll(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	_, err := s.Account().FindAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(a))

	accounts, err := s.Account().FindAll()
	assert.NoError(t, err)
	assert.NotNil(t, accounts)
}

func TestAccountRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	_, err := s.Account().FindByID(a.ID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(a))

	account, err := s.Account().FindByID(a.ID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
}

func TestAccountRepository_FindByTelegramID(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	_, err := s.Account().FindByTelegramID(a.TelegramID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(a))

	account, err := s.Account().FindByTelegramID(a.TelegramID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
}

func TestAccountRepository_FindBySuperuser(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL, migrations)
	defer teardown("account")

	s := sqlstore.New(db)
	a := model.TestAccount(t)

	_, err := s.Account().FindBySuperuser(a.Superuser)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	assert.NoError(t, s.Account().Create(a))

	accounts, err := s.Account().FindBySuperuser(a.Superuser)
	assert.NoError(t, err)
	assert.NotNil(t, accounts)
	assert.Equal(t, a.Superuser, accounts[0].Superuser)
}
