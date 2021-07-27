package store

import "gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"

type AccountRepository interface {
	Create(*model.Account) error
	FindByTelegramID(int64) (*model.Account, error)
}

type SchoolRepository interface {
	Create(*model.School) error
	FindByTitle(string) (*model.School, error)
	FindActive() (*model.School, error)
}

type LessonRepository interface {
	Create(*model.Lesson) error
	FindByID(int64) (*model.Lesson, error)
	FindByTitle(string) (*model.Lesson, error)
	FindBySchool(*model.School) ([]*model.Lesson, error)
}

type StudentRepository interface {
	Create(*model.Student) error
	FindByAccountSchool(*model.Account, *model.School) (*model.Student, error)
}

type HomeworkRepository interface {
	Create(*model.Homework) error
	FindByStudent(*model.Student) ([]*model.Homework, error)
	FindBySchool(*model.School) ([]*model.Homework, error)
	FindByStudentLesson(*model.Student, *model.Lesson) (*model.Homework, error)
}
