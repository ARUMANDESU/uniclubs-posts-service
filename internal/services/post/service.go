package postservice

import "errors"

var (
	ErrPostNotFound     = errors.New("post not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrClubNotFound     = errors.New("club not found")
	ErrInvalidArg       = errors.New("invalid argument")
)
