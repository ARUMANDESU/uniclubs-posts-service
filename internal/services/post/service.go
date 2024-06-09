package postservice

import (
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/client/club"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"log/slog"
)

var (
	ErrPostNotFound            = errors.New("post not found")
	ErrPermissionDenied        = errors.New("permission denied")
	ErrClubNotFound            = errors.New("club not found")
	ErrInvalidArg              = errors.New("invalid argument")
	ErrInvalidID               = errors.New("invalid id")
	ErrOptimisticLockingFailed = errors.New("optimistic locking failed")
)

func HandleError(log *slog.Logger, msg string, err error) error {
	switch {
	case errors.Is(err, club.ErrClubNotFound):
		return ErrClubNotFound
	case errors.Is(err, club.ErrInvalidArg):
		return ErrInvalidArg
	case errors.Is(err, storage.ErrNotFound):
		return ErrPostNotFound
	case errors.Is(err, storage.ErrInvalidID):
		return ErrInvalidID
	case errors.Is(err, storage.ErrOptimisticLockingFailed):
		return ErrOptimisticLockingFailed
	default:
		log.Error(msg, logger.Err(err))
		return err
	}
}
