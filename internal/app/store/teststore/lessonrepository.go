package teststore

import (
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
)

type LessonRepository struct {
	store   *Store
	lessons map[string]*model.Lesson
}

func (r *LessonRepository) Create(l *model.Lesson) error {
	if err := l.Validate(); err != nil {
		return err
	}

	r.lessons[l.Title] = l
	l.ID = int64(len(r.lessons))

	return nil
}

func (r *LessonRepository) FindByID(lesson_id int64) (*model.Lesson, error) {
	for _, l := range r.lessons {
		if l.ID == lesson_id {
			return l, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (r *LessonRepository) FindByTitle(title string) (*model.Lesson, error) {
	l, ok := r.lessons[title]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return l, nil
}

func (r *LessonRepository) FindBySchoolID(school_id int64) ([]*model.Lesson, error) {
	l := []*model.Lesson{}

	if r.store.studentRepository == nil {
		return nil, store.ErrRecordNotFound
	}

	students, ok := r.store.studentRepository.students[school_id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	if r.store.homeworkRepository == nil {
		return nil, store.ErrRecordNotFound
	}

	for _, student := range students {
		homeworks, ok := r.store.homeworkRepository.homeworks[student.ID]
		if !ok {
			continue
		}

		for _, homework := range homeworks {
			for _, lesson := range r.lessons {
				if lesson.ID == homework.Lesson.ID {
					l = appendLesson(l, lesson)
				}
			}
		}
	}

	if len(l) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return l, nil
}

func appendLesson(slice []*model.Lesson, homework *model.Lesson) []*model.Lesson {
	for _, elem := range slice {
		if elem.ID == homework.ID {
			return slice
		}
	}
	return append(slice, homework)
}
