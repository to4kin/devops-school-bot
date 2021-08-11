package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// HomeworkRepository ...
type HomeworkRepository struct {
	store     *Store
	homeworks []*model.Homework
}

// Create ...
func (r *HomeworkRepository) Create(h *model.Homework) error {
	if err := h.Validate(); err != nil {
		return err
	}

	homework, err := r.store.homeworkRepository.FindByStudentIDLessonID(h.Student.ID, h.Lesson.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if homework != nil {
		return store.ErrRecordIsExist
	}

	r.homeworks = append(r.homeworks, h)
	return nil
}

// FindByStudentID ...
func (r *HomeworkRepository) FindByStudentID(studentID int64) ([]*model.Homework, error) {
	hw := []*model.Homework{}

	for _, homework := range r.homeworks {
		if homework.Student.ID == studentID {
			hw = append(hw, homework)
		}
	}

	if len(hw) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return hw, nil
}

// FindBySchoolID ...
func (r *HomeworkRepository) FindBySchoolID(schoolID int64) ([]*model.Homework, error) {
	result := []*model.Homework{}

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

	for _, student := range students {
		for _, homework := range r.homeworks {
			if homework.Student.ID == student.ID {
				result = appendHomework(result, homework)
			}
		}
	}

	if len(result) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return result, nil
}

// FindByStudentIDLessonID ...
func (r *HomeworkRepository) FindByStudentIDLessonID(studentID int64, lessonID int64) (*model.Homework, error) {
	for _, homework := range r.homeworks {
		if homework.Student.ID == studentID && homework.Lesson.ID == lessonID {
			return homework, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func appendHomework(slice []*model.Homework, homework *model.Homework) []*model.Homework {
	for _, elem := range slice {
		if elem.ID == homework.ID {
			return slice
		}
	}
	return append(slice, homework)
}
