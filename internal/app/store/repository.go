package store

import "gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"

// AccountRepository ...
type AccountRepository interface {
	Create(*model.Account) error
	Update(*model.Account) error
	FindAll() ([]*model.Account, error)
	FindByID(int64) (*model.Account, error)
	FindByTelegramID(int64) (*model.Account, error)
}

// SchoolRepository ...
type SchoolRepository interface {
	Create(*model.School) error
	Update(*model.School) error
	FindAll() ([]*model.School, error)
	FindByID(int64) (*model.School, error)
	FindByTitle(string) (*model.School, error)
	FindByChatID(int64) (*model.School, error)
}

// LessonRepository ...
type LessonRepository interface {
	Create(*model.Lesson) error
	FindByID(int64) (*model.Lesson, error)
	FindByTitle(string) (*model.Lesson, error)
	FindBySchoolID(int64) ([]*model.Lesson, error)
}

// StudentRepository ...
type StudentRepository interface {
	Create(*model.Student) error
	FindAll() ([]*model.Student, error)
	FindByID(int64) (*model.Student, error)
	FindBySchoolID(int64) ([]*model.Student, error)
	FindByAccountIDSchoolID(int64, int64) (*model.Student, error)
}

// HomeworkRepository ...
type HomeworkRepository interface {
	Create(*model.Homework) error
	FindByStudentID(int64) ([]*model.Homework, error)
	FindBySchoolID(int64) ([]*model.Homework, error)
	FindByStudentIDLessonID(int64, int64) (*model.Homework, error)
}
