package service

import "errors"

var (
	ErrEmailAlreadyTaken = errors.New("email already taken")
)
