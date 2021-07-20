package store

import "gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"

type AccountRepository interface {
	Create(*model.Account) error
	FindByTelegramID(int64) (*model.Account, error)
}

type SchoolRepository interface {
	Create(*model.School) error
	FindByTitle(string) (*model.School, error)
}

type LessonRepository interface {
	Create(*model.Leson) error
	FindByTitle(string) (*model.Leson, error)
}
