package service

import "errors"

var (
	ErrEmailAlreadyTaken = errors.New("Email already taken")
	ErrInvalidCredentials = errors.New("Invalid email or password")
)
