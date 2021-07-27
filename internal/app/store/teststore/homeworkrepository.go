package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type HomeworkRepository struct {
	store     *Store
	homeworks map[int64]map[int64]*model.Homework
}

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

func (r *HomeworkRepository) FindByStudent(student *model.Student) ([]*model.Homework, error) {
	hw := []*model.Homework{}
	lessons, ok := r.homeworks[student.ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	for _, homework := range lessons {
		hw = append(hw, homework)
	}

	return hw, nil
}

func (r *HomeworkRepository) FindBySchool(school *model.School) ([]*model.Homework, error) {
	result := []*model.Homework{}

	if r.store.studentRepository == nil {
		return nil, store.ErrRecordNotFound
	}

	students, ok := r.store.studentRepository.students[school.ID]
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

func (r *HomeworkRepository) FindByStudentLesson(student *model.Student, lesson *model.Lesson) (*model.Homework, error) {
	lessons, ok := r.homeworks[student.ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	homework, ok := lessons[lesson.ID]
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
