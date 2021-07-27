package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type SchoolRepository struct {
	store   *Store
	schools map[string]*model.School
}

func (r *SchoolRepository) Create(s *model.School) error {
	if err := s.Validate(); err != nil {
		return err
	}

	r.schools[s.Title] = s
	s.ID = int64(len(r.schools))

	return nil
}

func (r *SchoolRepository) FindByTitle(title string) (*model.School, error) {
	s, ok := r.schools[title]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return s, nil
}

func (r *SchoolRepository) FindActive() (*model.School, error) {
	for _, school := range r.schools {
		if school.Active {
			return school, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
