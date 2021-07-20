package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type Store struct {
	userRepository     *UserRepository
	schoolRepository   *SchoolRepository
	homeworkRepository *HomeworkRepository
}

func New() *Store {
	return &Store{}
}

func (store *Store) User() store.UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}

	store.userRepository = &UserRepository{
		store: store,
		users: make(map[int64]*model.User),
	}

	return store.userRepository
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

func (store *Store) Homework() store.HomeworkRepository {
	if store.homeworkRepository != nil {
		return store.homeworkRepository
	}

	store.homeworkRepository = &HomeworkRepository{
		store:     store,
		homeworks: make(map[string]*model.Homework),
	}

	return store.homeworkRepository
}
