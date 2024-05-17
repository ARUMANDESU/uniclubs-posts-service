package storage

import "errors"

var (
	ErrUserExists              = errors.New("user already exists")
	ErrUserNotExists           = errors.New("user does not exist")
	ErrClubExists              = errors.New("club already exists")
	ErrClubNotExists           = errors.New("club does not exist")
	ErrEventNotFound           = errors.New("event not found")
	ErrOptimisticLockingFailed = errors.New("optimistic lock error")
	ErrInvalidID               = errors.New("the provided id is not a valid ObjectID")
	ErrInviteNotFound          = errors.New("invite not found")
)
