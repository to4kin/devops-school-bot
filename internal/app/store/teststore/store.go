package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// Store ...
type Store struct {
	accountRepository  *AccountRepository
	schoolRepository   *SchoolRepository
	lessonRepository   *LessonRepository
	moduleRepository   *ModuleRepository
	studentRepository  *StudentRepository
	homeworkRepository *HomeworkRepository
}

// New ...
func New() *Store {
	return &Store{}
}

// Account ...
func (store *Store) Account() store.AccountRepository {
	if store.accountRepository != nil {
		return store.accountRepository
	}

	store.accountRepository = &AccountRepository{
		store:    store,
		accounts: []*model.Account{},
	}

	return store.accountRepository
}

// School ...
func (store *Store) School() store.SchoolRepository {
	if store.schoolRepository != nil {
		return store.schoolRepository
	}

	store.schoolRepository = &SchoolRepository{
		store:   store,
		schools: []*model.School{},
	}

	return store.schoolRepository
}

// Lesson ...
func (store *Store) Lesson() store.LessonRepository {
	if store.lessonRepository != nil {
		return store.lessonRepository
	}

	store.lessonRepository = &LessonRepository{
		store:   store,
		lessons: []*model.Lesson{},
	}

	return store.lessonRepository
}

// Module ...
func (store *Store) Module() store.ModuleRepository {
	if store.moduleRepository != nil {
		return store.moduleRepository
	}

	store.moduleRepository = &ModuleRepository{
		store:   store,
		modules: []*model.Module{},
	}

	return store.moduleRepository
}

// Student ...
func (store *Store) Student() store.StudentRepository {
	if store.studentRepository != nil {
		return store.studentRepository
	}

	store.studentRepository = &StudentRepository{
		store:    store,
		students: []*model.Student{},
	}

	return store.studentRepository
}

// Homework ...
func (store *Store) Homework() store.HomeworkRepository {
	if store.homeworkRepository != nil {
		return store.homeworkRepository
	}

	store.homeworkRepository = &HomeworkRepository{
		store:     store,
		homeworks: []*model.Homework{},
	}

	return store.homeworkRepository
}
