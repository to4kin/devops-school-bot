package store

// Store ...
type Store interface {
	Account() AccountRepository
	School() SchoolRepository
	Lesson() LessonRepository
	Student() StudentRepository
	Homework() HomeworkRepository
}
