package store

import "gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByTelegramID(int64) (*model.User, error)
}

type SchoolRepository interface {
	Create(*model.School) error
	FindByTitle(string) (*model.School, error)
}
