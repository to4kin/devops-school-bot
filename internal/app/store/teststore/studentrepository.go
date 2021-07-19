package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type StudentRepository struct {
	store    *Store
	students map[int64]*model.Student
}

func (r *StudentRepository) Create(student *model.Student) error {
	if err := student.Validate(); err != nil {
		return err
	}

	r.students[student.TelegramID] = student
	student.ID = int64(len(r.students))

	return nil
}

func (r *StudentRepository) FindByTelegramID(telegramID int64) (*model.Student, error) {
	student, ok := r.students[telegramID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return student, nil
}
