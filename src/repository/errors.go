package repository

import "errors"

var (
	ErrTaskNotFound        = errors.New("task not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionAlreadyExist = errors.New("session already exist")
)
