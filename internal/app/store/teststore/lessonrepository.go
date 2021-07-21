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

func (r *LessonRepository) FindByTitle(title string) (*model.Lesson, error) {
	l, ok := r.lessons[title]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return l, nil
}
