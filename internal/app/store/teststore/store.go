package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type Store struct {
	accountRepository *AccountRepository
	schoolRepository  *SchoolRepository
	lessonRepository  *LessonRepository
}

func New() *Store {
	return &Store{}
}

func (store *Store) Account() store.AccountRepository {
	if store.accountRepository != nil {
		return store.accountRepository
	}

	store.accountRepository = &AccountRepository{
		store:    store,
		accounts: make(map[int64]*model.Account),
	}

	return store.accountRepository
}

func (store *Store) School() store.SchoolRepository {
	if store.schoolRepository != nil {
		return store.schoolRepository
	}

	store.schoolRepository = &SchoolRepository{
		store:   store,
		schools: make(map[string]*model.School),
	}

	return store.schoolRepository
}

func (store *Store) Lesson() store.LessonRepository {
	if store.lessonRepository != nil {
		return store.lessonRepository
	}

	store.lessonRepository = &LessonRepository{
		store:   store,
		lessons: make(map[string]*model.Leson),
	}

	return store.lessonRepository
}
