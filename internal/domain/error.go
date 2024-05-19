package domain

import "errors"

var (
	ErrOrganizerNotFound = errors.New("organizer not found")
	ErrUserIsEventOwner  = errors.New("user is event owner")
)
