package store

// Store ...
type Store interface {
	Account() AccountRepository
	School() SchoolRepository
	Lesson() LessonRepository
	Module() ModuleRepository
	Student() StudentRepository
	Homework() HomeworkRepository
	Callback() CallbackRepository
}
