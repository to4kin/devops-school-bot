package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	// ErrAnotherSchoolIsActive ...
	ErrAnotherSchoolIsActive = errors.New("another school is active")
)
