package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type Store struct {
	db                 *sql.DB
	userRepository     *UserRepository
	schoolRepository   *SchoolRepository
	homeworkRepository *HomeworkRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) User() store.UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}

	store.userRepository = &UserRepository{
		store: store,
	}

	return store.userRepository
}

func (store *Store) School() store.SchoolRepository {
	if store.schoolRepository != nil {
		return store.schoolRepository
	}

	store.schoolRepository = &SchoolRepository{
		store: store,
	}

	return store.schoolRepository
}

func (store *Store) Homework() store.HomeworkRepository {
	if store.homeworkRepository != nil {
		return store.homeworkRepository
	}

	store.homeworkRepository = &HomeworkRepository{
		store: store,
	}

	return store.homeworkRepository
}
