package usecases

import "errors"

var (
	ErrUserPassNotVerified = errors.New("password invalid")
	ErrUserExist           = errors.New("user already register")
)
