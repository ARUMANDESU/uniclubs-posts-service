package event

import "errors"

var (
	ErrClubNotExists           = errors.New("club not found")
	ErrEventNotFound           = errors.New("event not found")
	ErrEventUpdateConflict     = errors.New("event update conflict")
	ErrUserIsNotEventOwner     = errors.New("permissions denied: user is not event owner")
	ErrUserIsNotEventOrganizer = errors.New("user is not event organizer")
	ErrInvalidID               = errors.New("the provided id is not a valid ObjectID")
	ErrInviteAlreadyExists     = errors.New("invite already exists")
	ErrUserAlreadyOrganizer    = errors.New("user is already an organizer")
	ErrClubAlreadyCollaborator = errors.New("club is already a collaborator")
	ErrUserIsFromAnotherClub   = errors.New("user is from another club")
	ErrPermissionsDenied       = errors.New("permissions denied")
	ErrUserIsEventOwner        = errors.New("user is event owner")
	ErrInviteNotFound          = errors.New("invite not found")
)
