package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// ModuleRepository ...
type ModuleRepository struct {
	store   *Store
	modules []*model.Module
}

// Create ...
func (r *ModuleRepository) Create(m *model.Module) error {
	if err := m.Validate(); err != nil {
		return err
	}

	module, err := r.store.moduleRepository.FindByTitle(m.Title)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if module != nil {
		return store.ErrRecordIsExist
	}

	m.ID = int64(len(r.modules) + 1)
	r.modules = append(r.modules, m)
	return nil
}

// FindAll ...
func (r *ModuleRepository) FindAll() ([]*model.Module, error) {
	if len(r.modules) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return r.modules, nil
}

// FindByID ...
func (r *ModuleRepository) FindByID(moduleID int64) (*model.Module, error) {
	for _, m := range r.modules {
		if m.ID == moduleID {
			return m, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindByTitle ...
func (r *ModuleRepository) FindByTitle(title string) (*model.Module, error) {
	for _, module := range r.modules {
		if module.Title == title {
			return module, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindBySchoolID ...
func (r *ModuleRepository) FindBySchoolID(schoolID int64) ([]*model.Module, error) {
	m := []*model.Module{}

	if r.store.studentRepository == nil {
		return nil, store.ErrRecordNotFound
	}

	students := []*model.Student{}
	for _, student := range r.store.studentRepository.students {
		if student.School.ID == schoolID {
			students = append(students, student)
		}
	}

	if len(students) == 0 {
		return nil, store.ErrRecordNotFound
	}

	if r.store.homeworkRepository == nil {
		return nil, store.ErrRecordNotFound
	}

	for _, student := range students {
		for _, homework := range r.store.homeworkRepository.homeworks {
			if homework.Student.ID == student.ID {
				m = appendModule(m, homework.Lesson.Module)
			}
		}
	}

	if len(m) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return m, nil
}

func appendModule(slice []*model.Module, homework *model.Module) []*model.Module {
	for _, elem := range slice {
		if elem.ID == homework.ID {
			return slice
		}
	}
	return append(slice, homework)
}
