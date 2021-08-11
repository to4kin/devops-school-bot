package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	// ErrRecordIsExist ...
	ErrRecordIsExist = errors.New("record is already exist")
)
