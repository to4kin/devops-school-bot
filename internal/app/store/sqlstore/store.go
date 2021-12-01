package sqlstore

import (
	"database/sql"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// Store ...
type Store struct {
	db                 *sql.DB
	accountRepository  *AccountRepository
	schoolRepository   *SchoolRepository
	lessonRepository   *LessonRepository
	moduleRepository   *ModuleRepository
	studentRepository  *StudentRepository
	homeworkRepository *HomeworkRepository
	callbackRepository *CallbackRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Account ...
func (store *Store) Account() store.AccountRepository {
	if store.accountRepository != nil {
		return store.accountRepository
	}

	store.accountRepository = &AccountRepository{
		store: store,
	}

	return store.accountRepository
}

// School ...
func (store *Store) School() store.SchoolRepository {
	if store.schoolRepository != nil {
		return store.schoolRepository
	}

	store.schoolRepository = &SchoolRepository{
		store: store,
	}

	return store.schoolRepository
}

// Lesson ...
func (store *Store) Lesson() store.LessonRepository {
	if store.lessonRepository != nil {
		return store.lessonRepository
	}

	store.lessonRepository = &LessonRepository{
		store: store,
	}

	return store.lessonRepository
}

// Module ...
func (store *Store) Module() store.ModuleRepository {
	if store.moduleRepository != nil {
		return store.moduleRepository
	}

	store.moduleRepository = &ModuleRepository{
		store: store,
	}

	return store.moduleRepository
}

// Student ...
func (store *Store) Student() store.StudentRepository {
	if store.studentRepository != nil {
		return store.studentRepository
	}

	store.studentRepository = &StudentRepository{
		store: store,
	}

	return store.studentRepository
}

// Homework ...
func (store *Store) Homework() store.HomeworkRepository {
	if store.homeworkRepository != nil {
		return store.homeworkRepository
	}

	store.homeworkRepository = &HomeworkRepository{
		store: store,
	}

	return store.homeworkRepository
}

// Callback ...
func (store *Store) Callback() store.CallbackRepository {
	if store.callbackRepository != nil {
		return store.callbackRepository
	}

	store.callbackRepository = &CallbackRepository{
		store: store,
	}

	return store.callbackRepository
}
