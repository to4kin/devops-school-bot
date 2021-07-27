package store

import "errors"

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrAnotherSchoolIsActive = errors.New("another school is active")
)
