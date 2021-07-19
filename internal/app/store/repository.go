package store

import "gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"

type StudentRepository interface {
	Create(*model.Student) error
	FindByTelegramID(int64) (*model.Student, error)
}
