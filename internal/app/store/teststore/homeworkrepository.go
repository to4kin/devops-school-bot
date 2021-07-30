package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// HomeworkRepository ...
type HomeworkRepository struct {
	store     *Store
	homeworks map[int64]map[int64]*model.Homework
}

// Create ...
func (r *HomeworkRepository) Create(h *model.Homework) error {
	if err := h.Validate(); err != nil {
		return err
	}

	val := make(map[int64]*model.Homework)
	val[h.Lesson.ID] = h
	r.homeworks[h.Student.ID] = val
	h.ID = int64(len(r.homeworks))

	return nil
}

// FindByStudentID ...
func (r *HomeworkRepository) FindByStudentID(studentID int64) ([]*model.Homework, error) {
	hw := []*model.Homework{}
	lessons, ok := r.homeworks[studentID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	for _, homework := range lessons {
		hw = append(hw, homework)
	}

	return hw, nil
}

// FindBySchoolID ...
func (r *HomeworkRepository) FindBySchoolID(schoolID int64) ([]*model.Homework, error) {
	result := []*model.Homework{}

	if r.store.studentRepository == nil {
		return nil, store.ErrRecordNotFound
	}

	students, ok := r.store.studentRepository.students[schoolID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	for _, student := range students {
		homeworks, ok := r.homeworks[student.ID]
		if !ok {
			continue
		}

		for _, homework := range homeworks {
			result = appendHomework(result, homework)
		}
	}

	if len(result) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return result, nil
}

// FindByStudentIDLessonID ...
func (r *HomeworkRepository) FindByStudentIDLessonID(studentID int64, lessonID int64) (*model.Homework, error) {
	lessons, ok := r.homeworks[studentID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	homework, ok := lessons[lessonID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return homework, nil
}

func appendHomework(slice []*model.Homework, homework *model.Homework) []*model.Homework {
	for _, elem := range slice {
		if elem.ID == homework.ID {
			return slice
		}
	}
	return append(slice, homework)
}
