package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	// ErrSchoolIsExist ...
	ErrSchoolIsExist = errors.New("school is already exist")
)
