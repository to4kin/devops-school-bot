package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type Store struct {
	accountRepository  *AccountRepository
	schoolRepository   *SchoolRepository
	lessonRepository   *LessonRepository
	studentRepository  *StudentRepository
	homeworkRepository *HomeworkRepository
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
		lessons: make(map[string]*model.Lesson),
	}

	return store.lessonRepository
}

func (store *Store) Student() store.StudentRepository {
	if store.studentRepository != nil {
		return store.studentRepository
	}

	store.studentRepository = &StudentRepository{
		store:    store,
		students: make(map[int64]map[int64]*model.Student),
	}

	return store.studentRepository
}

func (store *Store) Homework() store.HomeworkRepository {
	if store.homeworkRepository != nil {
		return store.homeworkRepository
	}

	store.homeworkRepository = &HomeworkRepository{
		store:     store,
		homeworks: make(map[int64]map[int64]*model.Homework),
	}

	return store.homeworkRepository
}
