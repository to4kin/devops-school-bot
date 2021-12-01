package store

import "gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"

// AccountRepository ...
type AccountRepository interface {
	Create(*model.Account) error
	Update(*model.Account) error
	FindAll() ([]*model.Account, error)
	FindByID(int64) (*model.Account, error)
	FindByTelegramID(int64) (*model.Account, error)
	FindBySuperuser(bool) ([]*model.Account, error)
}

// SchoolRepository ...
type SchoolRepository interface {
	Create(*model.School) error
	Update(*model.School) error
	FindAll() ([]*model.School, error)
	FindByID(int64) (*model.School, error)
	FindByTitle(string) (*model.School, error)
	FindByChatID(int64) (*model.School, error)
	FindByActive(bool) ([]*model.School, error)
}

// LessonRepository ...
type LessonRepository interface {
	Create(*model.Lesson) error
	FindAll() ([]*model.Lesson, error)
	FindByID(int64) (*model.Lesson, error)
	FindByTitle(string) (*model.Lesson, error)
	FindBySchoolID(int64) ([]*model.Lesson, error)
}

// ModuleRepository ...
type ModuleRepository interface {
	Create(*model.Module) error
	FindAll() ([]*model.Module, error)
	FindByID(int64) (*model.Module, error)
	FindByTitle(string) (*model.Module, error)
	FindBySchoolID(int64) ([]*model.Module, error)
}

// StudentRepository ...
type StudentRepository interface {
	Create(*model.Student) error
	Update(*model.Student) error
	FindAll() ([]*model.Student, error)
	FindByID(int64) (*model.Student, error)
	FindBySchoolID(int64) ([]*model.Student, error)
	FindByAccountID(int64) ([]*model.Student, error)
	FindByAccountIDSchoolID(int64, int64) (*model.Student, error)
	FindByFullCourseSchoolID(bool, int64) ([]*model.Student, error)
}

// HomeworkRepository ...
type HomeworkRepository interface {
	Create(*model.Homework) error
	Update(*model.Homework) error
	DeleteByMessageID(int64) error
	FindByID(int64) (*model.Homework, error)
	FindByStudentID(int64) ([]*model.Homework, error)
	FindBySchoolID(int64) ([]*model.Homework, error)
	FindBySchoolIDLessonID(int64, int64) ([]*model.Homework, error)
	FindByStudentIDLessonID(int64, int64) (*model.Homework, error)
}

// CallbackRepository ...
type CallbackRepository interface {
	Create(*model.Callback) error
	FindByID(int64) (*model.Callback, error)
	FindByCallback(*model.Callback) (*model.Callback, error)
}
