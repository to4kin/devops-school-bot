package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type StudentRepository struct {
	store    *Store
	students map[int64]map[int64]*model.Student
}

func (r *StudentRepository) Create(s *model.Student) error {
	if err := s.Validate(); err != nil {
		return err
	}

	val := make(map[int64]*model.Student)
	val[s.School.ID] = s
	r.students[s.Account.ID] = val
	s.ID = int64(len(r.students))

	return nil
}

func (r *StudentRepository) FindByAccountSchool(account *model.Account, school *model.School) (*model.Student, error) {
	schools, ok := r.students[account.ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	student, ok := schools[school.ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return student, nil
}
