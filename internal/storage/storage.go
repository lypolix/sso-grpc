package storage

import "errors"


var (
	ErrUserExists = errors.New("user with this email already exists")
	ErrUserNotFound = errors.New("user with this email not found")
	ErrAppNotFound = errors.New("app with this id not found")
)