package store

type Store interface {
	Account() AccountRepository
	School() SchoolRepository
	Lesson() LessonRepository
}
