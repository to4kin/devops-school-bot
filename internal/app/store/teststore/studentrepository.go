package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// StudentRepository ...
type StudentRepository struct {
	store    *Store
	students map[int64]map[int64]*model.Student
}

// Create ...
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

// FindByAccountIDSchoolID ...
func (r *StudentRepository) FindByAccountIDSchoolID(accountID int64, schoolID int64) (*model.Student, error) {
	schools, ok := r.students[accountID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	student, ok := schools[schoolID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return student, nil
}
