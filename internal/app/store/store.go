package store

type Store interface {
	User() UserRepository
	School() SchoolRepository
}
