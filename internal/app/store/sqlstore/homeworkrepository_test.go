package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/sqlstore"
)

func TestHomeworkRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homeworks")

	s := sqlstore.New(db)
	h := model.TestHomework(t)

	assert.NoError(t, s.Homework().Create(h))
	assert.NotNil(t, h)
}

func TestHomework_FindByTitle(t *testing.T) {
	db, teardown := sqlstore.TestDb(t, databaseURL)
	defer teardown("homeworks")

	s := sqlstore.New(db)
	h := model.TestHomework(t)

	_, err := s.Homework().FindByTitle(h.Title)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Homework().Create(h)

	homework, err := s.Homework().FindByTitle(h.Title)
	assert.NoError(t, err)
	assert.NotNil(t, homework)
}
