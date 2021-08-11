package teststore

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// LessonRepository ...
type LessonRepository struct {
	store   *Store
	lessons []*model.Lesson
}

// Create ...
func (r *LessonRepository) Create(l *model.Lesson) error {
	if err := l.Validate(); err != nil {
		return err
	}

	lesson, err := r.store.lessonRepository.FindByTitle(l.Title)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	if lesson != nil {
		return store.ErrRecordIsExist
	}

	r.lessons = append(r.lessons, l)
	return nil
}

// FindByID ...
func (r *LessonRepository) FindByID(lessonID int64) (*model.Lesson, error) {
	for _, l := range r.lessons {
		if l.ID == lessonID {
			return l, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindByTitle ...
func (r *LessonRepository) FindByTitle(title string) (*model.Lesson, error) {
	for _, lesson := range r.lessons {
		if lesson.Title == title {
			return lesson, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

// FindBySchoolID ...
func (r *LessonRepository) FindBySchoolID(schoolID int64) ([]*model.Lesson, error) {
	l := []*model.Lesson{}

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
				l = appendLesson(l, homework.Lesson)
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
