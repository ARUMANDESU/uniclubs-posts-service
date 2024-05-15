package storage

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotExists = errors.New("user does not exist")
	ErrClubExists    = errors.New("club already exists")
	ErrClubNotExists = errors.New("club does not exist")
	ErrEventNotFound = errors.New("event not found")
)
