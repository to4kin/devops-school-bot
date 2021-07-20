package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/teststore"
)

func TestLesson_Create(t *testing.T) {
	s := teststore.New()
	l := model.TestLesson(t)

	assert.NoError(t, s.Lesson().Create(l))
	assert.NotNil(t, l)
}

func TestLesson_FindByTitle(t *testing.T) {
	s := teststore.New()
	l := model.TestLesson(t)

	_, err := s.Lesson().FindByTitle(l.Title)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Lesson().Create(l)

	lesson, err := s.Lesson().FindByTitle(l.Title)
	assert.NoError(t, err)
	assert.NotNil(t, lesson)
}
