package eventservice

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
	ErrUserIsFromAnotherClub   = errors.New("user is not member of the collaborator clubs")
	ErrPermissionsDenied       = errors.New("permissions denied")
	ErrUserIsEventOwner        = errors.New("user is event owner")
	ErrClubIsEventOwner        = errors.New("club is event owner")
	ErrInviteNotFound          = errors.New("invite not found")
	ErrCollaboratorNotFound    = errors.New("collaborator not found, club is not collaborator ")
	ErrOrganizerNotFound       = errors.New("organizer not found")
	ErrClubMismatch            = errors.New("club mismatch")
	ErrInvalidEventStatus      = errors.New("invalid event status")
	ErrEventInvalidFields      = errors.New("invalid event fields")
	ErrEventIsNotApproved      = errors.New("event is not approved")
	ErrEventIsNotEditable      = errors.New("event is not editable")
	ErrContainsUnchangeable    = errors.New("contains unchangeable fields")
	ErrUnknownStatus           = errors.New("unknown status")
	ErrEventIsFull             = errors.New("event is full")
	ErrAlreadyParticipating    = errors.New("user is already participating in the event")
	ErrParticipantNotFound     = errors.New("participant not found")
	ErrBanRecordNotFound       = errors.New("ban record not found")
	ErrUserAlreadyBanned       = errors.New("user is already banned")
	ErrUserIsBanned            = errors.New("user is banned")
)
