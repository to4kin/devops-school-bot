package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// SchoolRepository ...
type SchoolRepository struct {
	store   *Store
	schools []*model.School
}

// Create ...
func (r *SchoolRepository) Create(s *model.School) error {
	if err := s.Validate(); err != nil {
		return err
	}

	school, err := r.store.schoolRepository.FindByChatID(s.ChatID)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if school != nil {
		return store.ErrRecordIsExist
	}

	r.schools = append(r.schools, s)
	return nil
}

// ReActivate ...
func (r *SchoolRepository) ReActivate(s *model.School) error {
	for _, school := range r.schools {
		if school.ID == s.ID {
			s.Finished = false
			school = s
			return nil
		}
	}

	return store.ErrRecordNotFound
}

// Finish ...
func (r *SchoolRepository) Finish(s *model.School) error {
	for _, school := range r.schools {
		if school.ID == s.ID {
			s.Finished = true
			school = s
			return nil
		}
	}

	return store.ErrRecordNotFound
}

// FindAll ...
func (r *SchoolRepository) FindAll() ([]*model.School, error) {
	if len(r.schools) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return r.schools, nil
}

// FindByTitle ...
func (r *SchoolRepository) FindByTitle(title string) (*model.School, error) {
	for _, school := range r.schools {
		if school.Title == title {
			return school, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindByChatID ...
func (r *SchoolRepository) FindByChatID(chatID int64) (*model.School, error) {
	for _, school := range r.schools {
		if school.ChatID == chatID {
			return school, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
