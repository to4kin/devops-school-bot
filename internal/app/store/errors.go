package store

import "errors"

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrAnotherSchoolActive = errors.New("another school is active. please finish it and then start a new school")
)
