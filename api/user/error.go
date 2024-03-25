package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrWrongPassword      = errors.New("wrong password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrWrongRoute         = errors.New("cant use this route to update")
)
