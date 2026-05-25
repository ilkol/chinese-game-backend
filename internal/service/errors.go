package service

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidUserPassword = errors.New("invalid password")
)
