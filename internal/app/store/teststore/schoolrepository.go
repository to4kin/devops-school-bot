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

	s.ID = int64(len(r.schools) + 1)
	r.schools = append(r.schools, s)
	return nil
}

// Update ...
func (r *SchoolRepository) Update(s *model.School) error {
	for _, school := range r.schools {
		if school.ID == s.ID {
			school.Title = s.Title
			school.Active = s.Active
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

// FindByID ...
func (r *SchoolRepository) FindByID(id int64) (*model.School, error) {
	for _, school := range r.schools {
		if school.ID == id {
			return school, nil
		}
	}

	return nil, store.ErrRecordNotFound
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

// FindByActive ...
func (r *SchoolRepository) FindByActive(active bool) ([]*model.School, error) {
	schools := []*model.School{}
	for _, school := range r.schools {
		if school.Active == active {
			schools = append(schools, school)
		}
	}

	if len(schools) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return schools, nil
}
