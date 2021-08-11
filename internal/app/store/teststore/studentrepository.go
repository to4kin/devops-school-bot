package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// StudentRepository ...
type StudentRepository struct {
	store    *Store
	students []*model.Student
}

// Create ...
func (r *StudentRepository) Create(s *model.Student) error {
	if err := s.Validate(); err != nil {
		return err
	}

	student, err := r.store.studentRepository.FindByAccountIDSchoolID(s.Account.ID, s.School.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if student != nil {
		return store.ErrRecordIsExist
	}

	r.students = append(r.students, s)
	return nil
}

// FindBySchoolID ...
func (r *StudentRepository) FindBySchoolID(schoolID int64) ([]*model.Student, error) {
	result := []*model.Student{}

	for _, student := range r.students {
		if student.School.ID == schoolID {
			result = append(result, student)
		}
	}

	if len(result) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return result, nil
}

// FindByAccountIDSchoolID ...
func (r *StudentRepository) FindByAccountIDSchoolID(accountID int64, schoolID int64) (*model.Student, error) {
	for _, s := range r.students {
		if s.Account.ID == accountID && s.School.ID == schoolID {
			return s, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
