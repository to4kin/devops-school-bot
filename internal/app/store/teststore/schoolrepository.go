package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// SchoolRepository ...
type SchoolRepository struct {
	store   *Store
	schools map[string]*model.School
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
		return store.ErrSchoolIsExist
	}

	r.schools[s.Title] = s
	s.ID = int64(len(r.schools))

	return nil
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

// FindByTitle ...
func (r *SchoolRepository) FindByTitle(title string) (*model.School, error) {
	s, ok := r.schools[title]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return s, nil
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
