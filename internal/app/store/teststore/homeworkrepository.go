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
