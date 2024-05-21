package domain

import "errors"

var (
	ErrOrganizerNotFound    = errors.New("organizer not found")
	ErrUserIsEventOwner     = errors.New("user is event owner")
	ErrClubIsEventOwner     = errors.New("club is event owner")
	ErrCollaboratorsEmpty   = errors.New("there is no collaborators")
	ErrOrganizersEmpty      = errors.New("there is no organizers")
	ErrCollaboratorNotFound = errors.New("collaborator not found")
	ErrEventIsNotApproved   = errors.New("event is not approved")
)
