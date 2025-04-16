package repository

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrTaskNotFound        = ErrNotFound
	ErrUserNotFound        = ErrNotFound
	ErrSessionNotFound     = ErrNotFound
	ErrSessionAlreadyExist = errors.New("session already exist")
)
