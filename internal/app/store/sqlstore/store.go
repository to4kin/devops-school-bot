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
	studentRepository  *StudentRepository
	homeworkRepository *HomeworkRepository
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
