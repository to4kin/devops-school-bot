package store

import "gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"

type AccountRepository interface {
	Create(*model.Account) error
	FindByTelegramID(int64) (*model.Account, error)
}

type SchoolRepository interface {
	Create(*model.School) error
	FindByTitle(string) (*model.School, error)
}

type LessonRepository interface {
	Create(*model.Lesson) error
	FindByTitle(string) (*model.Lesson, error)
}

type StudentRepository interface {
	Create(*model.Student) error
	FindByAccountIDSchoolID(int64, int64) (*model.Student, error)
}

type HomeworkRepository interface {
	Create(*model.Homework) error
	FindByStudentIDLessonID(int64, int64) (*model.Homework, error)
}
